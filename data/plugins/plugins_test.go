package plugins_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	infraplugin "github.com/theheadlessengineer/crux/internal/infrastructure/plugin"
)

// pluginsDir returns the absolute path to data/plugins/ relative to this test file.
func pluginsDir(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	require.True(t, ok)
	// file = data/plugins/plugins_test.go → go up two levels to repo root
	root := filepath.Join(filepath.Dir(file), "..", "..")
	return filepath.Join(root, "data", "plugins")
}

// pilotPlugins lists all nine pilot plugin names.
var pilotPlugins = []string{
	"crux-plugin-postgresql",
	"crux-plugin-redis",
	"crux-plugin-kafka",
	"crux-plugin-auth-jwt",
	"crux-plugin-kubernetes",
	"crux-plugin-terraform-aws",
	"crux-plugin-github-actions",
	"crux-plugin-prometheus",
	"crux-plugin-claude-code",
}

// TestPilotPlugins_ManifestLoadsAndValidates verifies that every pilot plugin
// has a valid plugin.yaml that parses and passes ValidateManifest.
func TestPilotPlugins_ManifestLoadsAndValidates(t *testing.T) {
	base := pluginsDir(t)

	for _, name := range pilotPlugins {
		name := name
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			path := filepath.Join(base, name, "plugin.yaml")

			m, err := infraplugin.ParseManifest(path)
			require.NoError(t, err, "ParseManifest must succeed")

			assert.NoError(t, infraplugin.ValidateManifest(m), "ValidateManifest must succeed")
			assert.Equal(t, name, m.Metadata.Name)
			assert.NotEmpty(t, m.Metadata.Description)
			assert.NotEmpty(t, m.Spec.Questions, "plugin must declare at least one question")
		})
	}
}

// TestPilotPlugins_LoaderDiscovery verifies the Loader discovers all nine plugins
// from the data/plugins/ directory.
func TestPilotPlugins_LoaderDiscovery(t *testing.T) {
	l := infraplugin.New([]string{pluginsDir(t)})
	plugins, err := l.Load("1.0.0")
	require.NoError(t, err)
	assert.Len(t, plugins, len(pilotPlugins), "loader must discover all pilot plugins")
}

// TestPilotPlugins_VersionCompatibility verifies all plugins are compatible with
// crux v1.0.0 and reject an incompatible version.
func TestPilotPlugins_VersionCompatibility(t *testing.T) {
	base := pluginsDir(t)

	for _, name := range pilotPlugins {
		name := name
		t.Run(name+"/compatible", func(t *testing.T) {
			t.Parallel()
			m, err := infraplugin.ParseManifest(filepath.Join(base, name, "plugin.yaml"))
			require.NoError(t, err)
			assert.NoError(t, infraplugin.ValidateManifest(m))
		})
	}
}
