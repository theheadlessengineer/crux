package plugin_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domain "github.com/theheadlessengineer/crux/internal/domain/plugin"
	infraplugin "github.com/theheadlessengineer/crux/internal/infrastructure/plugin"
)

// validManifestYAML is a well-formed plugin.yaml.
const validManifestYAML = `
apiVersion: crux/v1
kind: Plugin
metadata:
  name: crux-plugin-postgresql
  version: 1.0.0
  description: PostgreSQL integration
  author: Platform Engineering
  trustTier: 1
  cruxVersionConstraint: ">=1.0.0"
spec:
  questions:
    - id: pg_version
      type: select
      prompt: "Which PostgreSQL version?"
      options: ["15", "16"]
      default: "16"
`

func writeManifest(t *testing.T, dir, content string) string {
	t.Helper()
	path := filepath.Join(dir, "plugin.yaml")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
	return path
}

// ── ParseManifest ─────────────────────────────────────────────────────────────

func TestParseManifest_ValidFile(t *testing.T) {
	dir := t.TempDir()
	path := writeManifest(t, dir, validManifestYAML)

	m, err := infraplugin.ParseManifest(path)
	require.NoError(t, err)
	assert.Equal(t, "crux-plugin-postgresql", m.Metadata.Name)
	assert.Equal(t, "1.0.0", m.Metadata.Version)
	assert.Equal(t, domain.TierOfficial, m.Metadata.TrustTier)
	assert.Equal(t, ">=1.0.0", m.Metadata.CruxVersionConstraint)
	assert.Len(t, m.Spec.Questions, 1)
	assert.Equal(t, "pg_version", m.Spec.Questions[0].ID)
}

func TestParseManifest_MissingFile(t *testing.T) {
	_, err := infraplugin.ParseManifest("/nonexistent/plugin.yaml")
	assert.Error(t, err)
}

func TestParseManifest_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	// Tabs in YAML mapping keys are invalid.
	path := writeManifest(t, dir, "apiVersion: crux/v1\n\t bad_indent: [unclosed")
	_, err := infraplugin.ParseManifest(path)
	assert.Error(t, err)
}

// ── ValidateManifest ──────────────────────────────────────────────────────────

func TestValidateManifest_Valid(t *testing.T) {
	dir := t.TempDir()
	m, err := infraplugin.ParseManifest(writeManifest(t, dir, validManifestYAML))
	require.NoError(t, err)
	assert.NoError(t, infraplugin.ValidateManifest(m))
}

func TestValidateManifest_MissingName(t *testing.T) {
	m := minimalManifest()
	m.Metadata.Name = ""
	assert.ErrorContains(t, infraplugin.ValidateManifest(m), "name")
}

func TestValidateManifest_BadVersion(t *testing.T) {
	m := minimalManifest()
	m.Metadata.Version = "not-semver"
	assert.ErrorContains(t, infraplugin.ValidateManifest(m), "semver")
}

func TestValidateManifest_BadAPIVersion(t *testing.T) {
	m := minimalManifest()
	m.APIVersion = "v2"
	assert.ErrorContains(t, infraplugin.ValidateManifest(m), "apiVersion")
}

func TestValidateManifest_BadKind(t *testing.T) {
	m := minimalManifest()
	m.Kind = "Service"
	assert.ErrorContains(t, infraplugin.ValidateManifest(m), "kind")
}

func TestValidateManifest_InvalidTrustTier(t *testing.T) {
	m := minimalManifest()
	m.Metadata.TrustTier = 99
	assert.ErrorContains(t, infraplugin.ValidateManifest(m), "trustTier")
}

func TestValidateManifest_MissingVersionConstraint(t *testing.T) {
	m := minimalManifest()
	m.Metadata.CruxVersionConstraint = ""
	assert.ErrorContains(t, infraplugin.ValidateManifest(m), "cruxVersionConstraint")
}

// ── Version compatibility ─────────────────────────────────────────────────────

func TestLoader_IncompatibleVersion_Rejected(t *testing.T) {
	const incompatible = `
apiVersion: crux/v1
kind: Plugin
metadata:
  name: crux-plugin-future
  version: 1.0.0
  description: Requires crux 99+
  author: Platform Engineering
  trustTier: 1
  cruxVersionConstraint: ">=99.0.0"
spec: {}
`
	dir := t.TempDir()
	pluginDir := filepath.Join(dir, "crux-plugin-future")
	require.NoError(t, os.MkdirAll(pluginDir, 0o755))
	writeManifest(t, pluginDir, incompatible)

	l := infraplugin.New([]string{dir})
	_, err := l.Load("1.0.0")
	assert.ErrorContains(t, err, "does not satisfy constraint")
}

func TestLoader_CompatibleVersion_Accepted(t *testing.T) {
	dir := t.TempDir()
	pluginDir := filepath.Join(dir, "crux-plugin-postgresql")
	require.NoError(t, os.MkdirAll(pluginDir, 0o755))
	writeManifest(t, pluginDir, validManifestYAML)

	l := infraplugin.New([]string{dir})
	plugins, err := l.Load("1.2.0")
	require.NoError(t, err)
	assert.Len(t, plugins, 1)
	assert.Equal(t, "crux-plugin-postgresql", plugins[0].Manifest.Metadata.Name)
}

// ── Hook execution ────────────────────────────────────────────────────────────

func TestRunPreGenerate_ExecutesInOrder(t *testing.T) {
	var order []string
	p := &domain.Plugin{
		Manifest: minimalManifest(),
		PreGenerate: []domain.Hook{
			func(_ context.Context, _ *domain.HookContext) error {
				order = append(order, "first")
				return nil
			},
			func(_ context.Context, _ *domain.HookContext) error {
				order = append(order, "second")
				return nil
			},
		},
	}
	require.NoError(t, infraplugin.RunPreGenerate(context.Background(), p, &domain.HookContext{}))
	assert.Equal(t, []string{"first", "second"}, order)
}

func TestRunPostGenerate_ExecutesInOrder(t *testing.T) {
	var order []string
	p := &domain.Plugin{
		Manifest: minimalManifest(),
		PostGenerate: []domain.Hook{
			func(_ context.Context, _ *domain.HookContext) error {
				order = append(order, "a")
				return nil
			},
			func(_ context.Context, _ *domain.HookContext) error {
				order = append(order, "b")
				return nil
			},
		},
	}
	require.NoError(t, infraplugin.RunPostGenerate(context.Background(), p, &domain.HookContext{}))
	assert.Equal(t, []string{"a", "b"}, order)
}

func TestRunPreGenerate_HookError_StopsExecution(t *testing.T) {
	called := false
	p := &domain.Plugin{
		Manifest: minimalManifest(),
		PreGenerate: []domain.Hook{
			func(_ context.Context, _ *domain.HookContext) error {
				return assert.AnError
			},
			func(_ context.Context, _ *domain.HookContext) error {
				called = true
				return nil
			},
		},
	}
	err := infraplugin.RunPreGenerate(context.Background(), p, &domain.HookContext{})
	assert.Error(t, err)
	assert.False(t, called, "second hook must not run after first fails")
}

// ── Loader discovery ──────────────────────────────────────────────────────────

func TestLoader_EmptyDir_ReturnsEmpty(t *testing.T) {
	l := infraplugin.New([]string{t.TempDir()})
	plugins, err := l.Load("1.0.0")
	require.NoError(t, err)
	assert.Empty(t, plugins)
}

func TestLoader_MultiplePlugins_LoadsAll(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"crux-plugin-redis", "crux-plugin-kafka"} {
		pd := filepath.Join(dir, name)
		require.NoError(t, os.MkdirAll(pd, 0o755))
		writeManifest(t, pd, minimalManifestYAML(name))
	}

	l := infraplugin.New([]string{dir})
	plugins, err := l.Load("1.0.0")
	require.NoError(t, err)
	assert.Len(t, plugins, 2)
}

func TestLoader_InvalidManifest_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	pd := filepath.Join(dir, "crux-plugin-bad")
	require.NoError(t, os.MkdirAll(pd, 0o755))
	// name is empty and version is not semver — both are validation errors.
	const badManifest = "apiVersion: crux/v1\nkind: Plugin\n" +
		"metadata:\n  name: \"\"\n  version: bad\n" +
		"  trustTier: 1\n  cruxVersionConstraint: \">=1.0.0\"\nspec: {}\n"
	writeManifest(t, pd, badManifest)

	l := infraplugin.New([]string{dir})
	_, err := l.Load("1.0.0")
	assert.Error(t, err)
}

// ── helpers ───────────────────────────────────────────────────────────────────

func minimalManifest() *domain.Manifest {
	return &domain.Manifest{
		APIVersion: "crux/v1",
		Kind:       "Plugin",
		Metadata: domain.Metadata{
			Name:                  "crux-plugin-test",
			Version:               "1.0.0",
			TrustTier:             domain.TierOfficial,
			CruxVersionConstraint: ">=1.0.0",
		},
	}
}

// minimalManifestYAML returns a valid plugin.yaml string for the given plugin name.
func minimalManifestYAML(name string) string {
	return "apiVersion: crux/v1\nkind: Plugin\n" +
		"metadata:\n  name: " + name + "\n  version: 1.0.0\n" +
		"  description: test\n  author: test\n" +
		"  trustTier: 1\n  cruxVersionConstraint: \">=1.0.0\"\nspec: {}\n"
}
