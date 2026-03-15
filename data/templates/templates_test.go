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
