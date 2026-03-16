package generator_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/theheadlessengineer/crux/internal/infrastructure/generator"
)

func cfg() generator.Config {
	return generator.Config{
		ServiceName: "payment-service",
		Module:      "github.com/company/payment-service",
		Language:    "go",
		Framework:   "gin",
		Team:        "payments",
		CLIVersion:  "1.0.0",
		GeneratedAt: time.Date(2026, 3, 16, 12, 0, 0, 0, time.UTC),
	}
}

func cfgPtr() *generator.Config {
	c := cfg()
	return &c
}

func TestGenerate_CreatesAllRequiredFiles(t *testing.T) {
	outDir := t.TempDir()
	require.NoError(t, generator.Generate(context.Background(), cfgPtr(), outDir))

	required := []string{
		"cmd/payment-service/main.go",
		"internal/config/config.go",
		"internal/presentation/http/router.go",
		"internal/presentation/http/health.go",
		"internal/presentation/http/server.go",
		"internal/infrastructure/logging/logger.go",
		"internal/infrastructure/errors/handler.go",
		"go.mod",
		"Makefile",
		"Dockerfile",
		".dockerignore",
		"docker-compose.yml",
		".gitignore",
		".editorconfig",
		".commitlintrc.yaml",
		".envrc",
		"README.md",
		"CHANGELOG.md",
		"resilience.yaml",
		"slo.yaml",
		"infra/kubernetes/deployment.yaml",
		"infra/kubernetes/networkpolicy-ingress.yaml",
		"infra/kubernetes/networkpolicy-egress.yaml",
		"infra/monitoring/alerts.yaml",
		"infra/monitoring/dashboard.json",
		"compliance/catalog-entry.yaml",
		"compliance/cost-budget.yaml",
		"compliance/data-classification.yaml",
		"compliance/log-retention.yaml",
		"docs/runbook.md",
		"docs/capacity-model.md",
		"docs/TODO.md",
		"docs/adr/ADR-001-initial-technology-choices.md",
		".github/workflows/ci.yaml",
		"scripts/seed.sh",
		"scripts/check_env.sh",
		"scripts/snapshot-db.sh",
		"scripts/restore-db.sh",
	}

	for _, rel := range required {
		assert.FileExists(t, filepath.Join(outDir, rel), "missing: %s", rel)
	}
}

func TestGenerate_CreatesEmptyDirStubs(t *testing.T) {
	outDir := t.TempDir()
	require.NoError(t, generator.Generate(context.Background(), cfgPtr(), outDir))

	stubs := []string{
		"internal/app/.gitkeep",
		"internal/domain/.gitkeep",
		"infra/terraform/.gitkeep",
		"tests/unit/.gitkeep",
		"tests/integration/.gitkeep",
	}
	for _, rel := range stubs {
		assert.FileExists(t, filepath.Join(outDir, rel), "missing stub: %s", rel)
	}
}

func TestGenerate_ScriptsAreExecutable(t *testing.T) {
	outDir := t.TempDir()
	require.NoError(t, generator.Generate(context.Background(), cfgPtr(), outDir))

	scripts := []string{
		"scripts/seed.sh",
		"scripts/check_env.sh",
		"scripts/snapshot-db.sh",
		"scripts/restore-db.sh",
	}
	for _, rel := range scripts {
		info, err := os.Stat(filepath.Join(outDir, rel))
		require.NoError(t, err, "stat %s", rel)
		assert.NotZero(t, info.Mode()&0o111, "%s must be executable", rel)
	}
}

func TestGenerate_RendersServiceName(t *testing.T) {
	outDir := t.TempDir()
	require.NoError(t, generator.Generate(context.Background(), cfgPtr(), outDir))

	mainGo, err := os.ReadFile(filepath.Join(outDir, "cmd", "payment-service", "main.go"))
	require.NoError(t, err)
	assert.Contains(t, string(mainGo), "payment-service")
}

func TestGenerate_RendersModulePath(t *testing.T) {
	outDir := t.TempDir()
	require.NoError(t, generator.Generate(context.Background(), cfgPtr(), outDir))

	goMod, err := os.ReadFile(filepath.Join(outDir, "go.mod"))
	require.NoError(t, err)
	assert.Contains(t, string(goMod), "github.com/company/payment-service")
}

func TestGenerate_DefaultModuleWhenEmpty(t *testing.T) {
	outDir := t.TempDir()
	c := cfg()
	c.Module = ""
	require.NoError(t, generator.Generate(context.Background(), &c, outDir))

	goMod, err := os.ReadFile(filepath.Join(outDir, "go.mod"))
	require.NoError(t, err)
	assert.Contains(t, string(goMod), "github.com/company/payment-service")
}

func TestGenerate_DockerfileNonRoot(t *testing.T) {
	outDir := t.TempDir()
	require.NoError(t, generator.Generate(context.Background(), cfgPtr(), outDir))

	dockerfile, err := os.ReadFile(filepath.Join(outDir, "Dockerfile"))
	require.NoError(t, err)
	assert.Contains(t, string(dockerfile), "nonroot")
}

func TestGenerate_ErrorOnInvalidOutputDir(t *testing.T) {
	err := generator.Generate(context.Background(), cfgPtr(), "/nonexistent/path/that/does/not/exist")
	assert.Error(t, err)
}
