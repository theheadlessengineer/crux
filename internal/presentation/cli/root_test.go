package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/theheadlessengineer/crux/internal/presentation/cli"
)

func TestRootCommand_Help(t *testing.T) {
	root := cli.BuildRoot("1.0.0", "abc123", "2026-01-01T00:00:00Z")
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"--help"})

	err := root.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "crux")
}

func TestRootCommand_Version(t *testing.T) {
	root := cli.BuildRoot("1.0.0", "abc123", "2026-01-01T00:00:00Z")
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"version"})

	err := root.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "1.0.0")
}

func TestRootCommand_GlobalFlags(t *testing.T) {
	root := cli.BuildRoot("1.0.0", "abc123", "2026-01-01T00:00:00Z")
	root.SetArgs([]string{"--help"})

	flags := root.PersistentFlags()
	assert.NotNil(t, flags.Lookup("verbose"))
	assert.NotNil(t, flags.Lookup("output"))
	assert.NotNil(t, flags.Lookup("config"))
}

func TestRootCommand_SubcommandsRegistered(t *testing.T) {
	root := cli.BuildRoot("1.0.0", "abc123", "2026-01-01T00:00:00Z")
	names := make([]string, 0, len(root.Commands()))
	for _, cmd := range root.Commands() {
		names = append(names, cmd.Name())
	}
	for _, want := range []string{"new", "version", "system", "validate"} {
		assert.True(t, contains(names, want), "expected subcommand %q to be registered", want)
	}
}

func contains(ss []string, s string) bool {
	for _, v := range ss {
		if strings.EqualFold(v, s) {
			return true
		}
	}
	return false
}
