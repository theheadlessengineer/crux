// Package generator wires the template engine and plugin system into the
// complete skeleton generation flow for crux new.
package generator

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	domain "github.com/theheadlessengineer/crux/internal/domain/template"
	infratemplate "github.com/theheadlessengineer/crux/internal/infrastructure/template"
)

// Config holds all inputs required to generate a service skeleton.
type Config struct {
	ServiceName string
	Module      string // Go module path, e.g. github.com/org/payment-service
	Language    string
	Framework   string
	Team        string
	CLIVersion  string
	GeneratedAt time.Time
}

// Generate renders all templates for the given config into outputDir.
func Generate(ctx context.Context, cfg *Config, outputDir string) error {
	eng, err := infratemplate.New()
	if err != nil {
		return fmt.Errorf("load templates: %w", err)
	}

	data := buildTemplateData(cfg)

	for tmplName, relPath := range fileMap(cfg.ServiceName) {
		outPath := filepath.Join(outputDir, relPath)
		if err := eng.Render(tmplName, data, outPath); err != nil {
			return fmt.Errorf("render %s: %w", tmplName, err)
		}
	}

	// Create empty-directory stubs (gitkeep files) for directories that have
	// no templates but must exist in the generated skeleton.
	for _, dir := range emptyDirs() {
		if err := mkGitkeep(filepath.Join(outputDir, dir)); err != nil {
			return err
		}
	}

	_ = ctx
	return nil
}

// fileMap returns the mapping of template name → relative output path.
// Template names match the embedded FS paths (e.g. "go-gin/cmd/main.go.tmpl").
func fileMap(serviceName string) map[string]string {
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
		// Root files
		"go-gin/go.mod.tmpl":             "go.mod",
		"go-gin/Makefile.tmpl":           "Makefile",
		"go-gin/Dockerfile.tmpl":         "Dockerfile",
		"go-gin/.dockerignore.tmpl":      ".dockerignore",
		"go-gin/docker-compose.yml.tmpl": "docker-compose.yml",
		"go-gin/.gitignore.tmpl":         ".gitignore",
		"go-gin/.editorconfig.tmpl":      ".editorconfig",
		"go-gin/.commitlintrc.yaml.tmpl": ".commitlintrc.yaml",
		"go-gin/.envrc.tmpl":             ".envrc",
		"go-gin/README.md.tmpl":          "README.md",
		"go-gin/CHANGELOG.md.tmpl":       "CHANGELOG.md",
		"go-gin/resilience.yaml.tmpl":    "resilience.yaml",
		"go-gin/slo.yaml.tmpl":           "slo.yaml",
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
		// CI
		"go-gin/github/workflows/ci.yaml.tmpl": ".github/workflows/ci.yaml",
		// Scripts
		"go-gin/scripts/seed.sh.tmpl":        "scripts/seed.sh",
		"go-gin/scripts/check_env.sh.tmpl":   "scripts/check_env.sh",
		"go-gin/scripts/snapshot-db.sh.tmpl": "scripts/snapshot-db.sh",
		"go-gin/scripts/restore-db.sh.tmpl":  "scripts/restore-db.sh",
	}
}

// emptyDirs lists directories that need a .gitkeep because no template writes into them.
func emptyDirs() []string {
	return []string{
		"internal/app",
		"internal/domain",
		"infra/terraform",
		"tests/unit",
		"tests/integration",
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
		Answers: map[string]any{},
	}
}
