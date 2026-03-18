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

// ── Multi-language tests ──────────────────────────────────────────────────────

func TestGenerate_PythonFastAPI_CoreFiles(t *testing.T) {
	outDir := t.TempDir()
	c := generator.Config{
		ServiceName: "order-service",
		Language:    "python",
		Framework:   "fastapi",
		Team:        "orders",
		CLIVersion:  "1.0.0",
		GeneratedAt: time.Date(2026, 3, 16, 12, 0, 0, 0, time.UTC),
	}
	require.NoError(t, generator.Generate(context.Background(), &c, outDir))

	required := []string{
		"main.py",
		"app/config.py",
		"app/health.py",
		"app/logging_config.py",
		"app/middleware.py",
		"requirements.txt",
		"Makefile",
		"Dockerfile",
		".github/workflows/ci.yaml",
		// Shared Tier 1 files
		"resilience.yaml",
		"slo.yaml",
		"infra/kubernetes/deployment.yaml",
		"infra/monitoring/alerts.yaml",
		"compliance/catalog-entry.yaml",
		"docs/runbook.md",
		"docs/TODO.md",
	}
	for _, rel := range required {
		assert.FileExists(t, filepath.Join(outDir, rel), "python: missing %s", rel)
	}
}

func TestGenerate_PythonFastAPI_DockerfileNonRoot(t *testing.T) {
	outDir := t.TempDir()
	c := generator.Config{ServiceName: "svc", Language: "python", CLIVersion: "1.0.0"}
	require.NoError(t, generator.Generate(context.Background(), &c, outDir))

	dockerfile, err := os.ReadFile(filepath.Join(outDir, "Dockerfile"))
	require.NoError(t, err)
	assert.Contains(t, string(dockerfile), "nonroot")
	assert.Contains(t, string(dockerfile), "USER nonroot")
}

func TestGenerate_JavaSpring_CoreFiles(t *testing.T) {
	outDir := t.TempDir()
	c := generator.Config{
		ServiceName: "user-service",
		Language:    "java",
		Framework:   "spring",
		Team:        "identity",
		CLIVersion:  "1.0.0",
		GeneratedAt: time.Date(2026, 3, 16, 12, 0, 0, 0, time.UTC),
	}
	require.NoError(t, generator.Generate(context.Background(), &c, outDir))

	required := []string{
		"src/main/java/Application.java",
		"src/main/java/health/HealthController.java",
		"src/main/resources/application.yaml",
		"pom.xml",
		"Makefile",
		"Dockerfile",
		".github/workflows/ci.yaml",
		// Shared Tier 1 files
		"resilience.yaml",
		"slo.yaml",
		"infra/kubernetes/deployment.yaml",
		"infra/monitoring/alerts.yaml",
		"compliance/catalog-entry.yaml",
		"docs/runbook.md",
		"docs/TODO.md",
	}
	for _, rel := range required {
		assert.FileExists(t, filepath.Join(outDir, rel), "java: missing %s", rel)
	}
}

func TestGenerate_JavaSpring_DockerfileNonRoot(t *testing.T) {
	outDir := t.TempDir()
	c := generator.Config{ServiceName: "svc", Language: "java", CLIVersion: "1.0.0"}
	require.NoError(t, generator.Generate(context.Background(), &c, outDir))

	dockerfile, err := os.ReadFile(filepath.Join(outDir, "Dockerfile"))
	require.NoError(t, err)
	assert.Contains(t, string(dockerfile), "nonroot")
}

func TestGenerate_NodeExpress_CoreFiles(t *testing.T) {
	outDir := t.TempDir()
	c := generator.Config{
		ServiceName: "notification-service",
		Language:    "node",
		Framework:   "express",
		Team:        "comms",
		CLIVersion:  "1.0.0",
		GeneratedAt: time.Date(2026, 3, 16, 12, 0, 0, 0, time.UTC),
	}
	require.NoError(t, generator.Generate(context.Background(), &c, outDir))

	required := []string{
		"index.js",
		"src/app.js",
		"src/health.js",
		"src/logging.js",
		"src/middleware.js",
		"package.json",
		"Makefile",
		"Dockerfile",
		".github/workflows/ci.yaml",
		// Shared Tier 1 files
		"resilience.yaml",
		"slo.yaml",
		"infra/kubernetes/deployment.yaml",
		"infra/monitoring/alerts.yaml",
		"compliance/catalog-entry.yaml",
		"docs/runbook.md",
		"docs/TODO.md",
	}
	for _, rel := range required {
		assert.FileExists(t, filepath.Join(outDir, rel), "node: missing %s", rel)
	}
}

func TestGenerate_NodeExpress_DockerfileNonRoot(t *testing.T) {
	outDir := t.TempDir()
	c := generator.Config{ServiceName: "svc", Language: "node", CLIVersion: "1.0.0"}
	require.NoError(t, generator.Generate(context.Background(), &c, outDir))

	dockerfile, err := os.ReadFile(filepath.Join(outDir, "Dockerfile"))
	require.NoError(t, err)
	assert.Contains(t, string(dockerfile), "nonroot")
}

func TestGenerate_AllLanguages_HealthEndpointsPresent(t *testing.T) {
	cases := []struct {
		language string
		file     string
		contains string
	}{
		{"go", "internal/presentation/http/health.go", "/health"},
		{"python", "app/health.py", "/health"},
		{"java", "src/main/java/health/HealthController.java", "/health"},
		{"node", "src/health.js", "/health"},
	}
	for _, tc := range cases {
		t.Run(tc.language, func(t *testing.T) {
			outDir := t.TempDir()
			c := generator.Config{ServiceName: "svc", Language: tc.language, CLIVersion: "1.0.0"}
			require.NoError(t, generator.Generate(context.Background(), &c, outDir))
			content, err := os.ReadFile(filepath.Join(outDir, tc.file))
			require.NoError(t, err)
			assert.Contains(t, string(content), tc.contains)
		})
	}
}

func TestGenerate_AllLanguages_SecurityHeadersPresent(t *testing.T) {
	cases := []struct {
		language string
		file     string
	}{
		{"go", "internal/presentation/http/router.go"},
		{"python", "app/middleware.py"},
		{"java", "src/main/resources/application.yaml"},
		{"node", "src/middleware.js"},
	}
	for _, tc := range cases {
		t.Run(tc.language, func(t *testing.T) {
			outDir := t.TempDir()
			c := generator.Config{ServiceName: "svc", Language: tc.language, CLIVersion: "1.0.0"}
			require.NoError(t, generator.Generate(context.Background(), &c, outDir))
			assert.FileExists(t, filepath.Join(outDir, tc.file))
		})
	}
}

// ── Plugin rendering tests ────────────────────────────────────────────────────

func TestGenerate_PluginSelected_RendersGoTemplates(t *testing.T) {
	outDir := t.TempDir()
	c := generator.Config{
		ServiceName: "payment-service",
		Language:    "go",
		CLIVersion:  "1.0.0",
		Plugins: []generator.SelectedPlugin{
			{
				Name: "crux-plugin-postgresql",
				Templates: []string{
					"internal/infrastructure/postgres/postgres.go.tmpl",
					"internal/infrastructure/postgres/health.go.tmpl",
				},
			},
		},
	}
	require.NoError(t, generator.Generate(context.Background(), &c, outDir))
	assert.FileExists(t, filepath.Join(outDir, "internal/infrastructure/postgres/postgres.go"))
	assert.FileExists(t, filepath.Join(outDir, "internal/infrastructure/postgres/health.go"))
}

func TestGenerate_PluginNotSelected_DoesNotRenderPluginFiles(t *testing.T) {
	outDir := t.TempDir()
	c := generator.Config{ServiceName: "payment-service", Language: "go", CLIVersion: "1.0.0"}
	require.NoError(t, generator.Generate(context.Background(), &c, outDir))
	assert.NoFileExists(t, filepath.Join(outDir, "internal/infrastructure/postgres/postgres.go"))
}

func TestGenerate_PluginSelected_JavaRendersJavaTemplate(t *testing.T) {
	outDir := t.TempDir()
	c := generator.Config{
		ServiceName: "user-service",
		Language:    "java",
		CLIVersion:  "1.0.0",
		Plugins: []generator.SelectedPlugin{
			{
				Name:      "crux-plugin-postgresql",
				Templates: []string{"src/main/java/infrastructure/postgres/PostgresConfig.java.tmpl"},
			},
		},
	}
	require.NoError(t, generator.Generate(context.Background(), &c, outDir))
	assert.FileExists(t, filepath.Join(outDir, "src/main/java/infrastructure/postgres/PostgresConfig.java"))
}

func TestGenerate_PluginSelected_PythonRendersPythonTemplate(t *testing.T) {
	outDir := t.TempDir()
	c := generator.Config{
		ServiceName: "order-service",
		Language:    "python",
		CLIVersion:  "1.0.0",
		Plugins: []generator.SelectedPlugin{
			{
				Name:      "crux-plugin-postgresql",
				Templates: []string{"app/db/postgres.py.tmpl"},
			},
		},
	}
	require.NoError(t, generator.Generate(context.Background(), &c, outDir))
	assert.FileExists(t, filepath.Join(outDir, "app/db/postgres.py"))
}

func TestGenerate_PluginSelected_NodeRendersNodeTemplate(t *testing.T) {
	outDir := t.TempDir()
	c := generator.Config{
		ServiceName: "notification-service",
		Language:    "node",
		CLIVersion:  "1.0.0",
		Plugins: []generator.SelectedPlugin{
			{
				Name:      "crux-plugin-postgresql",
				Templates: []string{"src/db/postgres.js.tmpl"},
			},
		},
	}
	require.NoError(t, generator.Generate(context.Background(), &c, outDir))
	assert.FileExists(t, filepath.Join(outDir, "src/db/postgres.js"))
}

func TestGenerate_MultiplePlugins_AllRendered(t *testing.T) {
	outDir := t.TempDir()
	c := generator.Config{
		ServiceName: "payment-service",
		Language:    "go",
		CLIVersion:  "1.0.0",
		Plugins: []generator.SelectedPlugin{
			{
				Name:      "crux-plugin-postgresql",
				Templates: []string{"internal/infrastructure/postgres/postgres.go.tmpl"},
			},
			{
				Name:      "crux-plugin-redis",
				Templates: []string{"internal/infrastructure/redis/redis.go.tmpl"},
			},
		},
	}
	require.NoError(t, generator.Generate(context.Background(), &c, outDir))
	assert.FileExists(t, filepath.Join(outDir, "internal/infrastructure/postgres/postgres.go"))
	assert.FileExists(t, filepath.Join(outDir, "internal/infrastructure/redis/redis.go"))
}
