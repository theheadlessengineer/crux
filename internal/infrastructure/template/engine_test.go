package template_test

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domain "github.com/theheadlessengineer/crux/internal/domain/template"
	infratmpl "github.com/theheadlessengineer/crux/internal/infrastructure/template"
)

// ── New / startup ─────────────────────────────────────────────────────────────

func TestNew_LoadsWithoutError(t *testing.T) {
	_, err := infratmpl.New()
	require.NoError(t, err, "New() must succeed when all embedded templates are valid")
}

func TestNew_InvalidTemplateSyntaxReturnsError(t *testing.T) {
	badFS := fstest.MapFS{
		"bad.tmpl": {Data: []byte(`{{ .Unclosed`)},
	}
	_, err := infratmpl.NewFromFS(badFS)
	assert.Error(t, err, "invalid template syntax must cause startup failure")
}

// ── Render ────────────────────────────────────────────────────────────────────

func TestRender_WritesRenderedFile(t *testing.T) {
	eng, err := infratmpl.New()
	require.NoError(t, err)

	out := filepath.Join(t.TempDir(), "Dockerfile")
	data := &domain.TemplateData{
		Service: domain.ServiceData{Name: "payment-service"},
	}
	require.NoError(t, eng.Render("go-gin/Dockerfile.tmpl", data, out))

	content, err := os.ReadFile(out)
	require.NoError(t, err)
	assert.Contains(t, string(content), "payment-service")
}

func TestRender_CreatesIntermediateDirectories(t *testing.T) {
	eng, err := infratmpl.New()
	require.NoError(t, err)

	out := filepath.Join(t.TempDir(), "deep", "nested", "Dockerfile")
	data := &domain.TemplateData{Service: domain.ServiceData{Name: "svc"}}

	require.NoError(t, eng.Render("go-gin/Dockerfile.tmpl", data, out))
	_, statErr := os.Stat(out)
	assert.NoError(t, statErr)
}

func TestRender_ShellScriptIsExecutable(t *testing.T) {
	eng, err := infratmpl.New()
	require.NoError(t, err)

	out := filepath.Join(t.TempDir(), "run.sh")
	data := &domain.TemplateData{Service: domain.ServiceData{Name: "svc"}}

	require.NoError(t, eng.Render("go-gin/Dockerfile.tmpl", data, out))

	info, err := os.Stat(out)
	require.NoError(t, err)
	assert.True(t, info.Mode()&0o111 != 0, "shell script must have executable bit set")
}

func TestRender_NonScriptIsNotExecutable(t *testing.T) {
	eng, err := infratmpl.New()
	require.NoError(t, err)

	out := filepath.Join(t.TempDir(), "Dockerfile")
	data := &domain.TemplateData{Service: domain.ServiceData{Name: "svc"}}

	require.NoError(t, eng.Render("go-gin/Dockerfile.tmpl", data, out))

	info, err := os.Stat(out)
	require.NoError(t, err)
	assert.True(t, info.Mode()&0o111 == 0, "non-script must not have executable bit set")
}

func TestRender_UnknownTemplateReturnsError(t *testing.T) {
	eng, err := infratmpl.New()
	require.NoError(t, err)

	err = eng.Render("does-not-exist.tmpl", &domain.TemplateData{}, filepath.Join(t.TempDir(), "out"))
	assert.ErrorContains(t, err, "not found")
}

func TestRender_VariableSubstitution(t *testing.T) {
	eng, err := infratmpl.New()
	require.NoError(t, err)

	out := filepath.Join(t.TempDir(), "resilience.yaml")
	data := &domain.TemplateData{
		Service: domain.ServiceData{Name: "order-service"},
	}
	require.NoError(t, eng.Render("go-gin/resilience.yaml.tmpl", data, out))

	content, err := os.ReadFile(out)
	require.NoError(t, err)
	assert.Contains(t, string(content), "order-service")
}

// ── Helper functions ──────────────────────────────────────────────────────────

func TestHelpers_StringFunctions(t *testing.T) {
	tests := []struct {
		name string
		tmpl string
		want string
	}{
		{"upper", `{{upper "hello"}}`, "HELLO"},
		{"lower", `{{lower "HELLO"}}`, "hello"},
		{"title", `{{title "hello world"}}`, "Hello World"},
		{"replace", `{{replace "hello world" "world" "go"}}`, "hello go"},
		{"camel_kebab", `{{camel "hello-world"}}`, "helloWorld"},
		{"camel_snake", `{{camel "hello_world"}}`, "helloWorld"},
		{"snake", `{{snake "hello-world"}}`, "hello_world"},
		{"kebab", `{{kebab "hello_world"}}`, "hello-world"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fsys := fstest.MapFS{
				"test.tmpl": {Data: []byte(tc.tmpl)},
			}
			eng, err := infratmpl.NewFromFS(fsys)
			require.NoError(t, err)

			out := filepath.Join(t.TempDir(), "out.txt")
			require.NoError(t, eng.Render("test.tmpl", &domain.TemplateData{}, out))

			content, err := os.ReadFile(out)
			require.NoError(t, err)
			assert.Equal(t, tc.want, string(content))
		})
	}
}

func TestHelpers_ContainsViaPluginsField(t *testing.T) {
	fsys := fstest.MapFS{
		"check.tmpl": {Data: []byte(`{{if contains .plugins_used "kafka"}}yes{{else}}no{{end}}`)},
	}
	eng, err := infratmpl.NewFromFS(fsys)
	require.NoError(t, err)

	tests := []struct {
		name    string
		plugins []string
		want    string
	}{
		{"found", []string{"postgres", "kafka"}, "yes"},
		{"not_found", []string{"postgres", "redis"}, "no"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := filepath.Join(t.TempDir(), "out.txt")
			data := &domain.TemplateData{Plugins: tc.plugins}
			require.NoError(t, eng.Render("check.tmpl", data, out))
			content, err := os.ReadFile(out)
			require.NoError(t, err)
			assert.Equal(t, tc.want, string(content))
		})
	}
}
