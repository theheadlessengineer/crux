package cli_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/theheadlessengineer/crux/internal/presentation/cli"
)

func TestVersionCommand_TextOutput(t *testing.T) {
	root := cli.BuildRoot("1.2.3", "deadbeef", "2026-03-15T00:00:00Z")
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"version"})

	require.NoError(t, root.Execute())
	out := buf.String()
	assert.Contains(t, out, "1.2.3")
	assert.Contains(t, out, "deadbeef")
	assert.Contains(t, out, "2026-03-15T00:00:00Z")
}

func TestVersionCommand_JSONOutput(t *testing.T) {
	root := cli.BuildRoot("1.2.3", "deadbeef", "2026-03-15T00:00:00Z")
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"--output", "json", "version"})

	require.NoError(t, root.Execute())

	var result map[string]string
	require.NoError(t, json.Unmarshal(buf.Bytes(), &result))
	assert.Equal(t, "1.2.3", result["version"])
	assert.Equal(t, "deadbeef", result["commit"])
	assert.Equal(t, "2026-03-15T00:00:00Z", result["buildTime"])
}

func TestSystemCommand_Runs(t *testing.T) {
	root := cli.BuildRoot("dev", "none", "unknown")
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"system"})

	// We don't assert pass/fail since CI may not have Docker.
	// We only assert the command runs and produces output.
	_ = root.Execute()
	assert.NotEmpty(t, buf.String())
}

func TestSystemCommand_JSONOutput(t *testing.T) {
	root := cli.BuildRoot("dev", "none", "unknown")
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"--output", "json", "system"})

	_ = root.Execute()

	var results []map[string]interface{}
	require.NoError(t, json.Unmarshal(buf.Bytes(), &results))
	assert.NotEmpty(t, results)
	for _, r := range results {
		assert.Contains(t, r, "name")
		assert.Contains(t, r, "status")
	}
}

func TestValidateCommand_MissingFiles(t *testing.T) {
	root := cli.BuildRoot("dev", "none", "unknown")
	var buf bytes.Buffer
	root.SetOut(&buf)
	// Point at a temp empty dir.
	root.SetArgs([]string{"validate", t.TempDir()})

	err := root.Execute()
	assert.Error(t, err)
	assert.Contains(t, buf.String(), "✘")
}

func TestValidateCommand_JSONOutput(t *testing.T) {
	root := cli.BuildRoot("dev", "none", "unknown")
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"--output", "json", "validate", t.TempDir()})

	_ = root.Execute()

	var results []map[string]interface{}
	require.NoError(t, json.Unmarshal(buf.Bytes(), &results))
	assert.NotEmpty(t, results)
}
