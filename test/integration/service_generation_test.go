//go:build integration

package integration

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/theheadlessengineer/crux/internal/infrastructure/generator"
)

// TestServiceGeneration_CompleteSkeleton generates a full service skeleton and
// verifies it compiles, all required files are present, and scripts are executable.
func TestServiceGeneration_CompleteSkeleton(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	outDir := t.TempDir()
	cfg := generator.Config{
		ServiceName: "test-service",
		Module:      "github.com/company/test-service",
		Language:    "go",
		Framework:   "gin",
		Team:        "platform",
		CLIVersion:  "1.0.0",
		GeneratedAt: time.Now().UTC(),
	}

	require.NoError(t, generator.Generate(context.Background(), &cfg, outDir))

	// Verify required files exist.
	required := []string{
		"cmd/test-service/main.go",
		"go.mod",
		"Makefile",
		"Dockerfile",
		"docker-compose.yml",
		".gitignore",
		".envrc",
		"resilience.yaml",
		"slo.yaml",
		"docs/runbook.md",
		"docs/capacity-model.md",
		"infra/kubernetes/deployment.yaml",
		"infra/monitoring/alerts.yaml",
		"scripts/seed.sh",
		"scripts/check_env.sh",
	}
	for _, rel := range required {
		assert.FileExists(t, filepath.Join(outDir, rel), "missing: %s", rel)
	}

	// Verify scripts are executable.
	for _, script := range []string{"scripts/seed.sh", "scripts/check_env.sh", "scripts/snapshot-db.sh", "scripts/restore-db.sh"} {
		info, err := os.Stat(filepath.Join(outDir, script))
		require.NoError(t, err)
		assert.NotZero(t, info.Mode()&0o111, "%s must be executable", script)
	}

	// Verify the generated service compiles.
	// go.mod only has stubs — run go mod tidy first, then build.
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = outDir
	cmd.Env = append(os.Environ(), "GOFLAGS=-mod=mod")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Logf("go mod tidy output: %s", out)
		// tidy may fail if dependencies aren't available in CI — skip compile check.
		t.Skip("go mod tidy failed (network unavailable?), skipping compile check")
	}

	build := exec.Command("go", "build", "./...")
	build.Dir = outDir
	out, err := build.CombinedOutput()
	if err != nil {
		t.Logf("go build output: %s", out)
	}
	assert.NoError(t, err, "generated service must compile")
}
