// Package template implements the template rendering engine.
package template

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"github.com/theheadlessengineer/crux/data/templates"
	domain "github.com/theheadlessengineer/crux/internal/domain/template"
)

// engine is the concrete implementation of domain.Engine.
type engine struct {
	tmpl *template.Template
}

// New loads all templates from the embedded filesystem and returns a ready Engine.
// Returns an error if any template has invalid syntax (Fail Fast).
func New() (domain.Engine, error) {
	return newFromFS(templates.FS)
}

// newFromFS builds an engine from any fs.FS — used by tests to inject inline templates.
func newFromFS(fsys fs.FS) (domain.Engine, error) {
	t := template.New("").Funcs(helperFuncs())

	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		if !strings.HasSuffix(path, ".tmpl") {
			return nil
		}
		raw, readErr := fs.ReadFile(fsys, path)
		if readErr != nil {
			return fmt.Errorf("read template %s: %w", path, readErr)
		}
		if _, parseErr := t.New(path).Parse(string(raw)); parseErr != nil {
			return fmt.Errorf("parse template %s: %w", path, parseErr)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("load templates: %w", err)
	}

	return &engine{tmpl: t}, nil
}

// Render executes the named template and writes the result to outputPath.
// Shell scripts (.sh) are written with executable permissions (0755).
func (e *engine) Render(templateName string, data *domain.TemplateData, outputPath string) error {
	t := e.tmpl.Lookup(templateName)
	if t == nil {
		return fmt.Errorf("template %q not found", templateName)
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0o750); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	perm := fs.FileMode(0o644)
	if strings.HasSuffix(outputPath, ".sh") {
		perm = 0o755
	}

	// outputPath is caller-controlled (not user input), so file inclusion is intentional.
	//nolint:gosec // G304: path is caller-supplied, not user input
	f, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	if err != nil {
		return fmt.Errorf("open output file: %w", err)
	}

	if err := t.Execute(f, data.ToMap()); err != nil {
		_ = f.Close()
		return fmt.Errorf("render template %q: %w", templateName, err)
	}
	return f.Close()
}

// helperFuncs returns the template function map required by the user story.
func helperFuncs() template.FuncMap {
	return template.FuncMap{
		"upper":    strings.ToUpper,
		"lower":    strings.ToLower,
		"title":    toTitle,
		"replace":  strings.ReplaceAll,
		"camel":    toCamel,
		"snake":    toSnake,
		"kebab":    toKebab,
		"contains": sliceContains,
	}
}

func toTitle(s string) string {
	if s == "" {
		return s
	}
	words := strings.Fields(s)
	for i, w := range words {
		if w == "" {
			continue
		}
		runes := []rune(w)
		runes[0] = unicode.ToUpper(runes[0])
		words[i] = string(runes)
	}
	return strings.Join(words, " ")
}

// toCamel converts a kebab-case or snake_case string to camelCase.
func toCamel(s string) string {
	parts := splitWords(s)
	for i, p := range parts {
		if i == 0 {
			parts[i] = strings.ToLower(p)
			continue
		}
		if p == "" {
			continue
		}
		runes := []rune(strings.ToLower(p))
		runes[0] = unicode.ToUpper(runes[0])
		parts[i] = string(runes)
	}
	return strings.Join(parts, "")
}

// toSnake converts a string to snake_case.
func toSnake(s string) string {
	return strings.Join(splitWords(s), "_")
}

// toKebab converts a string to kebab-case.
func toKebab(s string) string {
	return strings.Join(splitWords(s), "-")
}

// splitWords splits on hyphens, underscores, and spaces, lowercasing each part.
func splitWords(s string) []string {
	var parts []string
	for _, p := range strings.FieldsFunc(s, func(r rune) bool {
		return r == '-' || r == '_' || r == ' '
	}) {
		if p != "" {
			parts = append(parts, strings.ToLower(p))
		}
	}
	return parts
}

// sliceContains reports whether value is present in list.
func sliceContains(list []string, value string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}
