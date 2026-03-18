// Package generator wires the template engine and plugin system into the
// complete skeleton generation flow for crux new.
package generator

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	dataplugins "github.com/theheadlessengineer/crux/data/plugins"
	domain "github.com/theheadlessengineer/crux/internal/domain/template"
	infratemplate "github.com/theheadlessengineer/crux/internal/infrastructure/template"
)

// SelectedPlugin carries the name and language-resolved template list for a
// plugin the user selected during `crux new`.
type SelectedPlugin struct {
	Name      string
	Templates []string // resolved for the chosen language
}

// Config holds all inputs required to generate a service skeleton.
type Config struct {
	ServiceName string
	Module      string // Go module path, e.g. github.com/org/payment-service
	Language    string
	Framework   string
	Team        string
	CLIVersion  string
	GeneratedAt time.Time
	Plugins     []SelectedPlugin
	Answers     map[string]any
}

// Generate renders all templates for the given config into outputDir.
func Generate(ctx context.Context, cfg *Config, outputDir string) error {
	eng, err := infratemplate.New()
	if err != nil {
		return fmt.Errorf("load templates: %w", err)
	}

	data := buildTemplateData(cfg)

	// Render core language templates.
	for tmplName, relPath := range fileMap(cfg) {
		outPath := filepath.Join(outputDir, relPath)
		if err := eng.Render(tmplName, data, outPath); err != nil {
			return fmt.Errorf("render %s: %w", tmplName, err)
		}
	}

	// Render plugin templates.
	if err := renderPlugins(eng, cfg, data, outputDir); err != nil {
		return err
	}

	// Create empty-directory stubs (gitkeep files) for directories that have
	// no templates but must exist in the generated skeleton.
	for _, dir := range emptyDirs(cfg) {
		if err := mkGitkeep(filepath.Join(outputDir, dir)); err != nil {
			return err
		}
	}

	_ = ctx
	return nil
}

// renderPlugins loads each selected plugin's templates from the embedded FS
// and renders them into outputDir.
func renderPlugins(eng domain.Engine, cfg *Config, data *domain.TemplateData, outputDir string) error {
	if len(cfg.Plugins) == 0 {
		return nil
	}

	for _, sel := range cfg.Plugins {
		if len(sel.Templates) == 0 {
			continue
		}

		// Sub-FS rooted at the plugin's templates/ directory.
		pluginFS, err := fs.Sub(dataplugins.FS, sel.Name+"/templates")
		if err != nil {
			return fmt.Errorf("plugin %s: open templates fs: %w", sel.Name, err)
		}

		if err := eng.AddFromFS(pluginFS); err != nil {
			return fmt.Errorf("plugin %s: load templates: %w", sel.Name, err)
		}

		for _, tmplPath := range sel.Templates {
			outPath := filepath.Join(outputDir, strings.TrimSuffix(tmplPath, ".tmpl"))
			if err := eng.Render(tmplPath, data, outPath); err != nil {
				return fmt.Errorf("plugin %s: render %s: %w", sel.Name, tmplPath, err)
			}
		}
	}
	return nil
}

// fileMap returns the mapping of template name → relative output path for the
// selected language. Shared (language-agnostic) files are sourced from go-gin
// templates; language-specific files come from the language-specific directory.
func fileMap(cfg *Config) map[string]string {
	shared := sharedFileMap()

	switch cfg.Language {
	case "python", "python-fastapi":
		return mergeMaps(shared, pythonFastAPIFileMap())
	case "java", "java-spring":
		return mergeMaps(shared, javaSpringFileMap())
	case "node", "node-express":
		return mergeMaps(shared, nodeExpressFileMap())
	default: // "go", "go-gin", or empty
		return mergeMaps(shared, goGinFileMap(cfg.ServiceName))
	}
}

// sharedFileMap returns language-agnostic files (YAML, Markdown, CI stubs)
// that are identical across all languages. These are sourced from go-gin since
// they contain no Go-specific content.
func sharedFileMap() map[string]string {
	return map[string]string{
		// Operational config
		"go-gin/resilience.yaml.tmpl": "resilience.yaml",
		"go-gin/slo.yaml.tmpl":        "slo.yaml",
		// Kubernetes
		"go-gin/kubernetes/deployment.yaml.tmpl":            "infra/kubernetes/deployment.yaml",
		"go-gin/kubernetes/networkpolicy-ingress.yaml.tmpl": "infra/kubernetes/networkpolicy-ingress.yaml",
		"go-gin/kubernetes/networkpolicy-egress.yaml.tmpl":  "infra/kubernetes/networkpolicy-egress.yaml",
		// Monitoring
		"go-gin/monitoring/alerts.yaml.tmpl":    "infra/monitoring/alerts.yaml",
		"go-gin/monitoring/dashboard.json.tmpl": "infra/monitoring/dashboard.json",
		// Compliance
		"go-gin/compliance/catalog-entry.yaml.tmpl":       "compliance/catalog-entry.yaml",
		"go-gin/compliance/cost-budget.yaml.tmpl":         "compliance/cost-budget.yaml",
		"go-gin/compliance/data-classification.yaml.tmpl": "compliance/data-classification.yaml",
		"go-gin/compliance/log-retention.yaml.tmpl":       "compliance/log-retention.yaml",
		// Docs
		"go-gin/docs/runbook.md.tmpl":                                "docs/runbook.md",
		"go-gin/docs/capacity-model.md.tmpl":                         "docs/capacity-model.md",
		"go-gin/docs/TODO.md.tmpl":                                   "docs/TODO.md",
		"go-gin/docs/adr/ADR-001-initial-technology-choices.md.tmpl": "docs/adr/ADR-001-initial-technology-choices.md",
		// Root files shared across languages
		"go-gin/.editorconfig.tmpl":      ".editorconfig",
		"go-gin/.commitlintrc.yaml.tmpl": ".commitlintrc.yaml",
		"go-gin/CHANGELOG.md.tmpl":       "CHANGELOG.md",
		"go-gin/.gitignore.tmpl":         ".gitignore",
		"go-gin/.envrc.tmpl":             ".envrc",
		"go-gin/README.md.tmpl":          "README.md",
	}
}

func goGinFileMap(serviceName string) map[string]string {
	return map[string]string{
		// Application
		"go-gin/cmd/main.go.tmpl":                                   "cmd/" + serviceName + "/main.go",
		"go-gin/internal/config/config.go.tmpl":                     "internal/config/config.go",
		"go-gin/internal/domain/health/registry.go.tmpl":            "internal/domain/health/registry.go",
		"go-gin/internal/presentation/http/router.go.tmpl":          "internal/presentation/http/router.go",
		"go-gin/internal/presentation/http/health.go.tmpl":          "internal/presentation/http/health.go",
		"go-gin/internal/presentation/http/server.go.tmpl":          "internal/presentation/http/server.go",
		"go-gin/internal/infrastructure/logging/logger.go.tmpl":     "internal/infrastructure/logging/logger.go",
		"go-gin/internal/infrastructure/logging/middleware.go.tmpl": "internal/infrastructure/logging/middleware.go",
		"go-gin/internal/infrastructure/errors/handler.go.tmpl":     "internal/infrastructure/errors/handler.go",
		"go-gin/internal/infrastructure/tracing/provider.go.tmpl":   "internal/infrastructure/tracing/provider.go",
		"go-gin/internal/infrastructure/tracing/middleware.go.tmpl": "internal/infrastructure/tracing/middleware.go",
		"go-gin/internal/infrastructure/tracing/httpclient.go.tmpl": "internal/infrastructure/tracing/httpclient.go",
		"go-gin/internal/infrastructure/shutdown/shutdown.go.tmpl":  "internal/infrastructure/shutdown/shutdown.go",
		// Root
		"go-gin/go.mod.tmpl":             "go.mod",
		"go-gin/Makefile.tmpl":           "Makefile",
		"go-gin/Dockerfile.tmpl":         "Dockerfile",
		"go-gin/.dockerignore.tmpl":      ".dockerignore",
		"go-gin/docker-compose.yml.tmpl": "docker-compose.yml",
		// CI
		"go-gin/github/workflows/ci.yaml.tmpl": ".github/workflows/ci.yaml",
		// Scripts
		"go-gin/scripts/seed.sh.tmpl":        "scripts/seed.sh",
		"go-gin/scripts/check_env.sh.tmpl":   "scripts/check_env.sh",
		"go-gin/scripts/snapshot-db.sh.tmpl": "scripts/snapshot-db.sh",
		"go-gin/scripts/restore-db.sh.tmpl":  "scripts/restore-db.sh",
	}
}

func pythonFastAPIFileMap() map[string]string {
	return map[string]string{
		"python-fastapi/main.py.tmpl":                  "main.py",
		"python-fastapi/app/config.py.tmpl":            "app/config.py",
		"python-fastapi/app/health.py.tmpl":            "app/health.py",
		"python-fastapi/app/logging_config.py.tmpl":    "app/logging_config.py",
		"python-fastapi/app/middleware.py.tmpl":        "app/middleware.py",
		"python-fastapi/requirements.txt.tmpl":         "requirements.txt",
		"python-fastapi/Makefile.tmpl":                 "Makefile",
		"python-fastapi/Dockerfile.tmpl":               "Dockerfile",
		"python-fastapi/github/workflows/ci.yaml.tmpl": ".github/workflows/ci.yaml",
	}
}

func javaSpringFileMap() map[string]string {
	return map[string]string{
		"java-spring/src/main/java/Application.java.tmpl":             "src/main/java/Application.java",
		"java-spring/src/main/java/health/HealthController.java.tmpl": "src/main/java/health/HealthController.java",
		"java-spring/src/main/resources/application.yaml.tmpl":        "src/main/resources/application.yaml",
		"java-spring/pom.xml.tmpl":                                    "pom.xml",
		"java-spring/Makefile.tmpl":                                   "Makefile",
		"java-spring/Dockerfile.tmpl":                                 "Dockerfile",
		"java-spring/github/workflows/ci.yaml.tmpl":                   ".github/workflows/ci.yaml",
	}
}

func nodeExpressFileMap() map[string]string {
	return map[string]string{
		"node-express/index.js.tmpl":                 "index.js",
		"node-express/src/app.js.tmpl":               "src/app.js",
		"node-express/src/health.js.tmpl":            "src/health.js",
		"node-express/src/logging.js.tmpl":           "src/logging.js",
		"node-express/src/middleware.js.tmpl":        "src/middleware.js",
		"node-express/package.json.tmpl":             "package.json",
		"node-express/Makefile.tmpl":                 "Makefile",
		"node-express/Dockerfile.tmpl":               "Dockerfile",
		"node-express/github/workflows/ci.yaml.tmpl": ".github/workflows/ci.yaml",
	}
}

// emptyDirs lists directories that need a .gitkeep because no template writes into them.
func emptyDirs(cfg *Config) []string {
	switch cfg.Language {
	case "python", "python-fastapi":
		return []string{"tests/unit", "tests/integration", "infra/terraform"}
	case "java", "java-spring":
		return []string{"src/test/java", "infra/terraform"}
	case "node", "node-express":
		return []string{"tests", "infra/terraform"}
	default:
		return []string{
			"internal/app",
			"internal/domain",
			"infra/terraform",
			"tests/unit",
			"tests/integration",
		}
	}
}

func mkGitkeep(dir string) error {
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return fmt.Errorf("create directory %s: %w", dir, err)
	}
	path := filepath.Join(dir, ".gitkeep")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.WriteFile(path, nil, 0o600); err != nil {
			return fmt.Errorf("create .gitkeep in %s: %w", dir, err)
		}
	}
	return nil
}

func buildTemplateData(cfg *Config) *domain.TemplateData {
	module := cfg.Module
	if module == "" {
		module = "github.com/company/" + cfg.ServiceName
	}
	generatedAt := cfg.GeneratedAt
	if generatedAt.IsZero() {
		generatedAt = time.Now().UTC()
	}
	answers := cfg.Answers
	if answers == nil {
		answers = map[string]any{}
	}
	return &domain.TemplateData{
		Service: domain.ServiceData{
			Name:      cfg.ServiceName,
			Module:    module,
			Language:  cfg.Language,
			Framework: cfg.Framework,
			Team:      cfg.Team,
			Namespace: cfg.ServiceName,
		},
		Company: domain.CompanyData{
			Name:              "company",
			VaultAddr:         "https://vault.internal.company.com",
			CorrelationHeader: "traceparent",
		},
		Resilience: domain.ResilienceData{
			CircuitBreakerThreshold: 50,
			TimeoutDBMs:             3000,
			TimeoutHTTPMs:           5000,
			TimeoutKafkaMs:          10000,
			RetryMaxAttempts:        3,
			RetryBackoffBaseMs:      100,
		},
		SLO: domain.SLOData{
			AvailabilityTarget: "99.9",
			P99LatencyMs:       500,
			ErrorBudgetPolicy:  "halt_deployments_on_exhaustion",
		},
		Cost: domain.CostData{
			Centre:           "engineering",
			Team:             cfg.Team,
			MonthlyBudgetUSD: 500,
		},
		Infra: domain.InfraData{
			Cloud:  "aws",
			Region: "us-east-1",
		},
		Meta: domain.MetaData{
			CLIVersion:  cfg.CLIVersion,
			GeneratedAt: generatedAt.Format(time.RFC3339),
		},
		Plugins: []string{},
		Answers: answers,
	}
}

// mergeMaps merges src into dst (dst takes precedence on key collision).
func mergeMaps(dst, src map[string]string) map[string]string {
	result := make(map[string]string, len(dst)+len(src))
	for k, v := range dst {
		result[k] = v
	}
	for k, v := range src {
		result[k] = v
	}
	return result
}
