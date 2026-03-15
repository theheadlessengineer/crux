package cli_test

import (
	"bytes"
	"errors"
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
