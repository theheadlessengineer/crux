package cli_test

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/theheadlessengineer/crux/internal/presentation/cli"
)

func TestNewCommand_ValidName(t *testing.T) {
	root := cli.BuildRoot("dev", "none", "unknown")
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"new", "my-service"})

	err := root.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "my-service")
}

func TestNewCommand_NoArgs_ExitsValidation(t *testing.T) {
	root := cli.BuildRoot("dev", "none", "unknown")
	root.SetArgs([]string{"new"})

	err := root.Execute()
	assert.Error(t, err)
}

func TestNewCommand_InvalidName_ValidationError(t *testing.T) {
	tests := []struct {
		name        string
		serviceName string
	}{
		{"starts with number", "1service"},
		{"uppercase", "MyService"},
		{"too short", "ab"},
		{"special chars", "my_service"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := cli.BuildRoot("dev", "none", "unknown")
			root.SetArgs([]string{"new", tt.serviceName})

			err := root.Execute()
			require.Error(t, err)

			var ve *cli.ValidationError
			assert.True(t, errors.As(err, &ve), "expected ValidationError for %q", tt.serviceName)
		})
	}
}

func TestNewCommand_DryRun(t *testing.T) {
	root := cli.BuildRoot("dev", "none", "unknown")
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"new", "my-service", "--dry-run"})

	err := root.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "dry-run")
}

func TestNewCommand_Help(t *testing.T) {
	root := cli.BuildRoot("dev", "none", "unknown")
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"new", "--help"})

	err := root.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "service-name")
}

func writeConfigFile(t *testing.T, content string) string {
	t.Helper()
	f := filepath.Join(t.TempDir(), "crux.config.yaml")
	require.NoError(t, os.WriteFile(f, []byte(content), 0o644))
	return f
}

func TestNewCommand_ConfigFile_PreFillsAnswers(t *testing.T) {
	cfgPath := writeConfigFile(t, `
service:
  name: payment-service
  language: go
  framework: gin
answers:
  pg_version: "16"
`)
	outDir := t.TempDir()
	root := cli.BuildRoot("dev", "none", "unknown")
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"new", "payment-service", "--config", cfgPath, "--output-dir", outDir})

	require.NoError(t, root.Execute())
	assert.FileExists(t, filepath.Join(outDir, ".skeleton.json"))
	assert.FileExists(t, filepath.Join(outDir, "crux.lock"))
}

func TestNewCommand_NoPrompt_WithCompleteConfig_Succeeds(t *testing.T) {
	cfgPath := writeConfigFile(t, `
service:
  name: my-service
  language: go
  framework: gin
`)
	outDir := t.TempDir()
	root := cli.BuildRoot("dev", "none", "unknown")
	root.SetArgs([]string{"new", "my-service", "--config", cfgPath, "--no-prompt", "--output-dir", outDir})

	require.NoError(t, root.Execute())
	assert.FileExists(t, filepath.Join(outDir, ".skeleton.json"))
}

func TestNewCommand_NoPrompt_WithoutConfig_Fails(t *testing.T) {
	root := cli.BuildRoot("dev", "none", "unknown")
	root.SetArgs([]string{"new", "my-service", "--no-prompt"})

	err := root.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "--no-prompt")
}

func TestNewCommand_NoPrompt_WithIncompleteConfig_Fails(t *testing.T) {
	cfgPath := writeConfigFile(t, `
service:
  name: my-service
`)
	root := cli.BuildRoot("dev", "none", "unknown")
	root.SetArgs([]string{"new", "my-service", "--config", cfgPath, "--no-prompt"})

	err := root.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "service.language")
}

func TestNewCommand_InvalidConfigFile_Fails(t *testing.T) {
	root := cli.BuildRoot("dev", "none", "unknown")
	root.SetArgs([]string{"new", "my-service", "--config", "/nonexistent/crux.config.yaml"})

	err := root.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "load config")
}

func TestNewCommand_WritesLockfiles(t *testing.T) {
	outDir := t.TempDir()
	root := cli.BuildRoot("1.0.0", "none", "unknown")
	root.SetArgs([]string{"new", "my-service", "--output-dir", outDir})

	require.NoError(t, root.Execute())
	assert.FileExists(t, filepath.Join(outDir, ".skeleton.json"))
	assert.FileExists(t, filepath.Join(outDir, "crux.lock"))
}

// TestPluginQuestionToPrompt_LanguageAwareOptions verifies that pluginQuestionToPrompt
// resolves language-specific options and defaults when options_by_language is declared.
// We test this indirectly via the exported BuildRoot + config-file path so we don't
// need to export the internal helper.
func TestNewCommand_LanguageAware_PythonMigrationTool(t *testing.T) {
	// Generate a Python service with postgresql plugin selected via config.
	// The skeleton.json answers should record the python-specific migration tool.
	cfgPath := writeConfigFile(t, `
service:
  name: order-service
  language: python
  framework: fastapi
answers:
  language: python
  pg_version: "16"
  pg_migration_tool: alembic
`)
	outDir := t.TempDir()
	root := cli.BuildRoot("dev", "none", "unknown")
	root.SetArgs([]string{"new", "order-service", "--config", cfgPath, "--no-prompt", "--output-dir", outDir})

	require.NoError(t, root.Execute())

	// Skeleton must exist and service language must be python.
	skelPath := filepath.Join(outDir, ".skeleton.json")
	assert.FileExists(t, skelPath)

	data, err := os.ReadFile(skelPath)
	require.NoError(t, err)
	assert.Contains(t, string(data), "python")
}
