package templates_test

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func renderDockerfile(t *testing.T, data any) string {
	t.Helper()
	raw, err := os.ReadFile("go-gin/Dockerfile.tmpl")
	require.NoError(t, err, "Dockerfile.tmpl must exist")

	tmpl, err := template.New("Dockerfile").Parse(string(raw))
	require.NoError(t, err, "Dockerfile.tmpl must be valid Go template syntax")

	var buf bytes.Buffer
	require.NoError(t, tmpl.Execute(&buf, data))
	return buf.String()
}

func TestDockerfile_RendersWithServiceName(t *testing.T) {
	out := renderDockerfile(t, map[string]any{
		"service": map[string]any{"name": "payment-service"},
	})
	assert.Contains(t, out, "./cmd/payment-service", "binary path must use service name")
}

func TestDockerfile_MultiStageBuild(t *testing.T) {
	out := renderDockerfile(t, map[string]any{
		"service": map[string]any{"name": "svc"},
	})
	assert.Contains(t, out, "AS builder", "must have a named builder stage")
	assert.Contains(t, out, "FROM gcr.io/distroless/static:nonroot", "final stage must use distroless nonroot")
}

func TestDockerfile_NonRootUser(t *testing.T) {
	out := renderDockerfile(t, map[string]any{
		"service": map[string]any{"name": "svc"},
	})
	assert.Contains(t, out, "USER nonroot:nonroot", "process must run as nonroot user")
}

func TestDockerfile_CGODisabled(t *testing.T) {
	out := renderDockerfile(t, map[string]any{
		"service": map[string]any{"name": "svc"},
	})
	assert.Contains(t, out, "CGO_ENABLED=0", "CGO must be disabled for distroless compatibility")
}

func TestDockerfile_VersionBuildArg(t *testing.T) {
	out := renderDockerfile(t, map[string]any{
		"service": map[string]any{"name": "svc"},
	})
	assert.Contains(t, out, "ARG VERSION=dev", "VERSION build arg must be present with default")
	assert.Contains(t, out, "-X main.version=${VERSION}", "VERSION must be injected via ldflags")
}

func TestDockerfile_Healthcheck(t *testing.T) {
	out := renderDockerfile(t, map[string]any{
		"service": map[string]any{"name": "svc"},
	})
	assert.Contains(t, out, "HEALTHCHECK", "HEALTHCHECK instruction must be present")
}

func TestDockerfile_ExposesPort8080(t *testing.T) {
	out := renderDockerfile(t, map[string]any{
		"service": map[string]any{"name": "svc"},
	})
	assert.Contains(t, out, "EXPOSE 8080")
}

func TestDockerignore_ExcludesRequiredPaths(t *testing.T) {
	raw, err := os.ReadFile("go-gin/.dockerignore.tmpl")
	require.NoError(t, err)
	content := string(raw)

	for _, entry := range []string{".git/", "*.md", "docs/", "bin/"} {
		assert.True(t, strings.Contains(content, entry), ".dockerignore must exclude %q", entry)
	}
}

// ── Kubernetes manifest tests ─────────────────────────────────────────────────

var k8sData = map[string]any{
	"service": map[string]any{
		"name":        "payment-service",
		"namespace":   "payments",
		"environment": "production",
		"version":     "1.0.0",
	},
	"cost": map[string]any{
		"team": "payments",
	},
	"infra": map[string]any{
		"registry": "123456789.dkr.ecr.us-east-1.amazonaws.com",
	},
}

func renderK8s(t *testing.T, path string) string {
	t.Helper()
	raw, err := os.ReadFile(path)
	require.NoError(t, err, "%s must exist", path)
	tmpl, err := template.New(path).Parse(string(raw))
	require.NoError(t, err, "%s must be valid Go template syntax", path)
	var buf bytes.Buffer
	require.NoError(t, tmpl.Execute(&buf, k8sData))
	return buf.String()
}

func TestDeployment_ReadOnlyRootFilesystem(t *testing.T) {
	out := renderK8s(t, "go-gin/kubernetes/deployment.yaml.tmpl")
	assert.Contains(t, out, "readOnlyRootFilesystem: true")
}

func TestDeployment_RunsAsNonRoot(t *testing.T) {
	out := renderK8s(t, "go-gin/kubernetes/deployment.yaml.tmpl")
	assert.Contains(t, out, "runAsNonRoot: true")
	assert.Contains(t, out, "runAsUser: 65534")
}

func TestDeployment_AllowPrivilegeEscalationFalse(t *testing.T) {
	out := renderK8s(t, "go-gin/kubernetes/deployment.yaml.tmpl")
	assert.Contains(t, out, "allowPrivilegeEscalation: false")
}

func TestDeployment_DropsAllCapabilities(t *testing.T) {
	out := renderK8s(t, "go-gin/kubernetes/deployment.yaml.tmpl")
	assert.Contains(t, out, "drop:")
	assert.Contains(t, out, "- ALL")
}

func TestDeployment_TmpEmptyDirVolume(t *testing.T) {
	out := renderK8s(t, "go-gin/kubernetes/deployment.yaml.tmpl")
	assert.Contains(t, out, "mountPath: /tmp")
	assert.Contains(t, out, "emptyDir: {}")
}

func TestDeployment_RendersServiceName(t *testing.T) {
	out := renderK8s(t, "go-gin/kubernetes/deployment.yaml.tmpl")
	assert.Contains(t, out, "name: payment-service")
	assert.Contains(t, out, "namespace: payments")
}

func TestNetworkPolicyIngress_DefaultDeny(t *testing.T) {
	out := renderK8s(t, "go-gin/kubernetes/networkpolicy-ingress.yaml.tmpl")
	assert.Contains(t, out, "policyTypes:")
	assert.Contains(t, out, "- Ingress")
	assert.Contains(t, out, "ingress: []")
}

func TestNetworkPolicyEgress_DefaultDeny(t *testing.T) {
	out := renderK8s(t, "go-gin/kubernetes/networkpolicy-egress.yaml.tmpl")
	assert.Contains(t, out, "policyTypes:")
	assert.Contains(t, out, "- Egress")
}

func TestNetworkPolicyEgress_DNSAllowed(t *testing.T) {
	out := renderK8s(t, "go-gin/kubernetes/networkpolicy-egress.yaml.tmpl")
	assert.Contains(t, out, "port: 53")
}

func TestNetworkPolicies_RendersServiceName(t *testing.T) {
	for _, path := range []string{
		"go-gin/kubernetes/networkpolicy-ingress.yaml.tmpl",
		"go-gin/kubernetes/networkpolicy-egress.yaml.tmpl",
	} {
		out := renderK8s(t, path)
		assert.Contains(t, out, "payment-service", "service name must appear in %s", path)
		assert.Contains(t, out, "namespace: payments", "namespace must appear in %s", path)
	}
}

// ── Resilience template tests ─────────────────────────────────────────────────

var resilienceData = map[string]any{
	"service": map[string]any{"name": "payment-service"},
}

func TestResilienceYAML_ContainsRequiredSections(t *testing.T) {
	raw, err := os.ReadFile("go-gin/resilience.yaml.tmpl")
	require.NoError(t, err)
	tmpl, err := template.New("resilience").Parse(string(raw))
	require.NoError(t, err)
	var buf bytes.Buffer
	require.NoError(t, tmpl.Execute(&buf, resilienceData))
	out := buf.String()

	for _, section := range []string{"timeout:", "retry:", "circuitBreaker:", "bulkhead:", "mesh_mode:"} {
		assert.Contains(t, out, section, "resilience.yaml must contain %q", section)
	}
	assert.Contains(t, out, "payment-service")
}

// ── SLO template tests ────────────────────────────────────────────────────────

func TestSLOYAML_ContainsRequiredFields(t *testing.T) {
	raw, err := os.ReadFile("go-gin/slo.yaml.tmpl")
	require.NoError(t, err)
	tmpl, err := template.New("slo").Parse(string(raw))
	require.NoError(t, err)
	var buf bytes.Buffer
	require.NoError(t, tmpl.Execute(&buf, resilienceData))
	out := buf.String()

	for _, field := range []string{"availability", "latency-p99", "target:", "window:", "error_budget_policy:"} {
		assert.Contains(t, out, field, "slo.yaml must contain %q", field)
	}
}

// ── Alerts template tests ─────────────────────────────────────────────────────

func TestAlertsYAML_ContainsFourGoldenSignals(t *testing.T) {
	raw, err := os.ReadFile("go-gin/monitoring/alerts.yaml.tmpl")
	require.NoError(t, err)
	tmpl, err := template.New("alerts").Parse(string(raw))
	require.NoError(t, err)
	var buf bytes.Buffer
	require.NoError(t, tmpl.Execute(&buf, resilienceData))
	out := buf.String()

	for _, alert := range []string{"HighErrorRate", "HighP99Latency", "NoTraffic", "HighMemorySaturation"} {
		assert.Contains(t, out, alert, "alerts.yaml must contain alert %q", alert)
	}
	assert.Contains(t, out, "payment-service")
}

// ── Grafana dashboard template tests ─────────────────────────────────────────

func TestDashboardJSON_ContainsFourPanels(t *testing.T) {
	raw, err := os.ReadFile("go-gin/monitoring/dashboard.json.tmpl")
	require.NoError(t, err)
	tmpl, err := template.New("dashboard").Parse(string(raw))
	require.NoError(t, err)
	var buf bytes.Buffer
	require.NoError(t, tmpl.Execute(&buf, resilienceData))
	out := buf.String()

	for _, panel := range []string{"Request Rate", "Error Rate", "Latency", "Saturation"} {
		assert.Contains(t, out, panel, "dashboard.json must contain panel %q", panel)
	}
	assert.Contains(t, out, "payment-service")
}

// ── Compliance stub template tests ────────────────────────────────────────────

var complianceData = map[string]any{
	"service": map[string]any{
		"name":        "payment-service",
		"environment": "production",
	},
	"cost": map[string]any{
		"team":   "payments",
		"centre": "engineering",
	},
}

func TestCostBudgetYAML_ContainsRequiredFields(t *testing.T) {
	raw, err := os.ReadFile("go-gin/compliance/cost-budget.yaml.tmpl")
	require.NoError(t, err)
	tmpl, err := template.New("cost-budget").Parse(string(raw))
	require.NoError(t, err)
	var buf bytes.Buffer
	require.NoError(t, tmpl.Execute(&buf, complianceData))
	out := buf.String()

	for _, field := range []string{"monthly_budget_usd:", "alert_threshold_percent:", "components:"} {
		assert.Contains(t, out, field)
	}
	assert.Contains(t, out, "payment-service")
	assert.Contains(t, out, "payments")
}

func TestDataClassificationYAML_ContainsRequiredFields(t *testing.T) {
	raw, err := os.ReadFile("go-gin/compliance/data-classification.yaml.tmpl")
	require.NoError(t, err)
	tmpl, err := template.New("data-classification").Parse(string(raw))
	require.NoError(t, err)
	var buf bytes.Buffer
	require.NoError(t, tmpl.Execute(&buf, complianceData))
	out := buf.String()

	assert.Contains(t, out, "fields:")
	assert.Contains(t, out, "payment-service")
}

func TestLogRetentionYAML_ContainsAllEnvironments(t *testing.T) {
	raw, err := os.ReadFile("go-gin/compliance/log-retention.yaml.tmpl")
	require.NoError(t, err)
	tmpl, err := template.New("log-retention").Parse(string(raw))
	require.NoError(t, err)
	var buf bytes.Buffer
	require.NoError(t, tmpl.Execute(&buf, complianceData))
	out := buf.String()

	for _, env := range []string{"production:", "staging:", "development:"} {
		assert.Contains(t, out, env)
	}
	assert.Contains(t, out, "logs_days:")
}

func TestCatalogEntryYAML_ContainsRequiredFields(t *testing.T) {
	raw, err := os.ReadFile("go-gin/compliance/catalog-entry.yaml.tmpl")
	require.NoError(t, err)
	tmpl, err := template.New("catalog-entry").Parse(string(raw))
	require.NoError(t, err)
	var buf bytes.Buffer
	require.NoError(t, tmpl.Execute(&buf, complianceData))
	out := buf.String()

	for _, field := range []string{"ownership:", "lifecycle:", "dependencies:", "tags:"} {
		assert.Contains(t, out, field)
	}
	assert.Contains(t, out, "payment-service")
}

// ── CI workflow template tests ────────────────────────────────────────────────

func TestCIWorkflow_ContainsComplianceSteps(t *testing.T) {
	raw, err := os.ReadFile("go-gin/github/workflows/ci.yaml.tmpl")
	require.NoError(t, err)
	tmpl, err := template.New("ci").Parse(string(raw))
	require.NoError(t, err)
	var buf bytes.Buffer
	require.NoError(t, tmpl.Execute(&buf, resilienceData))
	out := buf.String()

	assert.Contains(t, out, "Generate SBOM", "CI must include SBOM generation step")
	assert.Contains(t, out, "Licence scan", "CI must include licence scan step")
	assert.Contains(t, out, "if: false", "DAST step must be disabled by default")
	assert.Contains(t, out, "DAST", "DAST slot must be present")
	assert.Contains(t, out, "pre-commit", "CI must run pre-commit hooks")
}

// ── Documentation stub template tests ────────────────────────────────────────

var docsData = map[string]any{
	"service": map[string]any{
		"name":         "payment-service",
		"language":     "go",
		"framework":    "gin",
		"service_type": "REST API",
		"namespace":    "payments",
	},
	"cost": map[string]any{
		"team": "payments",
	},
	"meta": map[string]any{
		"generated_at": "2026-03-15T17:00:00Z",
		"cli_version":  "1.0.0",
	},
	"plugins_used": []string{"crux-plugin-kubernetes@2.0.0"},
}

func TestADR001_ContainsChoices(t *testing.T) {
	raw, err := os.ReadFile("go-gin/docs/adr/ADR-001-initial-technology-choices.md.tmpl")
	require.NoError(t, err)
	tmpl, err := template.New("adr001").Parse(string(raw))
	require.NoError(t, err)
	var buf bytes.Buffer
	require.NoError(t, tmpl.Execute(&buf, docsData))
	out := buf.String()

	assert.Contains(t, out, "payment-service")
	assert.Contains(t, out, "go")
	assert.Contains(t, out, "gin")
	assert.Contains(t, out, "crux-plugin-kubernetes@2.0.0")
}

func TestRunbook_ContainsIncidentAndDRSections(t *testing.T) {
	raw, err := os.ReadFile("go-gin/docs/runbook.md.tmpl")
	require.NoError(t, err)
	tmpl, err := template.New("runbook").Parse(string(raw))
	require.NoError(t, err)
	var buf bytes.Buffer
	require.NoError(t, tmpl.Execute(&buf, docsData))
	out := buf.String()

	assert.Contains(t, out, "Incident Response")
	assert.Contains(t, out, "Disaster Recovery")
	assert.Contains(t, out, "payment-service")
}

func TestTODO_ContainsRequiredPlaceholders(t *testing.T) {
	raw, err := os.ReadFile("go-gin/docs/TODO.md.tmpl")
	require.NoError(t, err)
	tmpl, err := template.New("todo").Parse(string(raw))
	require.NoError(t, err)
	var buf bytes.Buffer
	require.NoError(t, tmpl.Execute(&buf, docsData))
	out := buf.String()

	required := []string{
		"cost-budget.yaml", "data-classification.yaml",
		"catalog-entry.yaml", "slo.yaml", "resilience.yaml",
	}
	for _, item := range required {
		assert.Contains(t, out, item, "TODO.md must reference %q", item)
	}
}

func TestEditorconfig_EnforcesUTF8AndLF(t *testing.T) {
	raw, err := os.ReadFile("go-gin/.editorconfig.tmpl")
	require.NoError(t, err)
	content := string(raw)
	assert.Contains(t, content, "charset = utf-8")
	assert.Contains(t, content, "end_of_line = lf")
	assert.Contains(t, content, "indent_style = tab")
}

func TestCommitlintrc_EnforcesConventionalCommits(t *testing.T) {
	raw, err := os.ReadFile("go-gin/.commitlintrc.yaml.tmpl")
	require.NoError(t, err)
	content := string(raw)
	assert.Contains(t, content, "config-conventional")
	assert.Contains(t, content, "feat")
	assert.Contains(t, content, "fix")
}
