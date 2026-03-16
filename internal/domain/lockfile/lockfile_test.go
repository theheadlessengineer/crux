package lockfile_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/theheadlessengineer/crux/internal/domain/lockfile"
)

func skeleton() *lockfile.Skeleton {
	return &lockfile.Skeleton{
		CruxVersion: "1.0.0",
		GeneratedAt: time.Date(2026, 3, 16, 12, 0, 0, 0, time.UTC),
		Service: lockfile.SkeletonService{
			Name:      "payment-service",
			Language:  "go",
			Framework: "gin",
		},
		Answers: map[string]any{"pg_version": "16"},
		Plugins: []lockfile.PluginEntry{
			{Name: "crux-plugin-postgresql", Version: "1.2.0"},
			{Name: "crux-plugin-redis", Version: "0.9.1"},
		},
		Deviations: []string{},
		Tier1Standards: lockfile.Tier1Standards{
			Enforced:          true,
			DisabledStandards: []string{},
		},
	}
}

func TestWrite_CreatesSkeletonJSON(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, lockfile.Write(dir, skeleton()))

	data, err := os.ReadFile(filepath.Join(dir, ".skeleton.json"))
	require.NoError(t, err)

	var got lockfile.Skeleton
	require.NoError(t, json.Unmarshal(data, &got))

	assert.Equal(t, "1.0.0", got.CruxVersion)
	assert.Equal(t, "payment-service", got.Service.Name)
	assert.Equal(t, "go", got.Service.Language)
	assert.Equal(t, "gin", got.Service.Framework)
	assert.Equal(t, "16", got.Answers["pg_version"])
	assert.Len(t, got.Plugins, 2)
	assert.True(t, got.Tier1Standards.Enforced)
	assert.Equal(t, time.Date(2026, 3, 16, 12, 0, 0, 0, time.UTC), got.GeneratedAt)
}

func TestWrite_CreatesLockfile(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, lockfile.Write(dir, skeleton()))

	data, err := os.ReadFile(filepath.Join(dir, "crux.lock"))
	require.NoError(t, err)

	var got lockfile.Lockfile
	require.NoError(t, json.Unmarshal(data, &got))

	assert.Equal(t, 1, got.LockfileVersion)
	assert.Equal(t, "1.2.0", got.Plugins["crux-plugin-postgresql"])
	assert.Equal(t, "0.9.1", got.Plugins["crux-plugin-redis"])
}

func TestWrite_LockfileContainsOnlyVersionPins(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, lockfile.Write(dir, skeleton()))

	data, err := os.ReadFile(filepath.Join(dir, "crux.lock"))
	require.NoError(t, err)

	// Lockfile must only have lockfileVersion and plugins — no other fields.
	var raw map[string]any
	require.NoError(t, json.Unmarshal(data, &raw))
	assert.ElementsMatch(t, []string{"lockfileVersion", "plugins"}, keys(raw))
}

func TestWrite_ValidJSON(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, lockfile.Write(dir, skeleton()))

	for _, name := range []string{".skeleton.json", "crux.lock"} {
		data, err := os.ReadFile(filepath.Join(dir, name))
		require.NoError(t, err)
		assert.True(t, json.Valid(data), "%s is not valid JSON", name)
	}
}

func TestWrite_EmptyPlugins(t *testing.T) {
	dir := t.TempDir()
	s := skeleton()
	s.Plugins = nil
	require.NoError(t, lockfile.Write(dir, s))

	data, err := os.ReadFile(filepath.Join(dir, "crux.lock"))
	require.NoError(t, err)

	var got lockfile.Lockfile
	require.NoError(t, json.Unmarshal(data, &got))
	assert.Empty(t, got.Plugins)
}

func TestWrite_ErrorOnInvalidDir(t *testing.T) {
	err := lockfile.Write("/nonexistent/path/that/does/not/exist", skeleton())
	assert.Error(t, err)
}

func keys(m map[string]any) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
