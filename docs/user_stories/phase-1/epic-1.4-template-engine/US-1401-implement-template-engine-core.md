# US-1401 — Implement Template Engine Core

**Epic:** 1.4 Template Engine
**Phase:** 1 — Pilot
**Priority:** Must Have
**Status:** Done

---

## User Story

As a platform engineer building Crux,
I want a template rendering engine that loads Go templates from embedded files and renders them to the filesystem with variable substitution,
so that every feature that generates files can use a consistent, testable rendering system.

---

## Pre-Development Checklist

- [x] Epic 0.1 Foundation is complete
- [x] The template variable namespace is agreed (struct fields available in templates)
- [x] The embed strategy is agreed (`//go:embed` directive)
- [x] Story estimated and accepted into the sprint

---

## Scope

Implement the template engine as an infrastructure layer component that wraps Go's `text/template` package with a well-defined interface, variable namespace, helper functions, and filesystem output.

### In Scope

- A `TemplateEngine` interface with `Render(templateName string, data TemplateData, outputPath string) error`
- Go `text/template` wrapper with the agreed helper function set
- Template variable namespace struct (`TemplateData`) covering: service name, language, framework, plugins, team, environment, version, and all answers from the prompt engine
- Template loading from embedded files (`data/templates/`)
- Rendering output to the target filesystem path
- File permission setting on rendered files (executable bit for shell scripts)
- Template validation at startup — invalid templates cause startup failure (Fail Fast)
- Unit tests for rendering, helper functions, and error cases

### Helper Functions Required

| Function | Description |
|---|---|
| `upper` | Convert string to uppercase |
| `lower` | Convert string to lowercase |
| `title` | Convert string to title case |
| `replace` | Replace substring |
| `camel` | Convert to camelCase |
| `snake` | Convert to snake_case |
| `kebab` | Convert to kebab-case |
| `contains` | Check if a list contains a value |

### Out of Scope

- Actual templates for specific languages (Epic 1.5)
- Plugin template loading (Epic 1.6)

---

## Acceptance Criteria

- [x] `TemplateEngine` interface is defined in the domain layer
  - `internal/domain/template/template.go` — `Engine` interface with `Render(*TemplateData, string) error`
- [x] Templates load from embedded files at startup
  - `data/templates/embed.go` — `//go:embed go-gin` exposes `FS embed.FS`; `New()` walks and parses all `.tmpl` files at construction time
- [x] Invalid template syntax causes a startup failure with a clear error
  - `newFromFS` returns a wrapped parse error; `TestNew_InvalidTemplateSyntaxReturnsError` confirms this
- [x] Variable substitution works for all `TemplateData` fields
  - `TemplateData.ToMap()` maps all fields to the lowercase key convention used by `.tmpl` files; `TestRender_VariableSubstitution` confirms substitution
- [x] All helper functions are implemented and tested
  - `helperFuncs()` in `engine.go` registers all 8 functions; `TestHelpers_StringFunctions` (8 sub-tests) and `TestHelpers_ContainsViaPluginsField` (2 sub-tests) cover every function
- [x] Rendered files have correct permissions (executable scripts, read-write for config)
  - `.sh` outputs get `0755`; all others get `0644`; `TestRender_ShellScriptIsExecutable` and `TestRender_NonScriptIsNotExecutable` verify this
- [x] Unit tests cover rendering, helper functions, and invalid template handling
  - 11 tests in `internal/infrastructure/template/engine_test.go`; all pass
- [x] `make test` passes

---

## Post-Completion Checklist

- [ ] Code reviewed by at least one other platform engineer
- [x] Each helper function tested manually
- [x] Template loading from embedded files confirmed at startup
- [x] Unit tests pass
- [ ] Story moved to Done in the project tracker

---

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| Epic 0.1 Foundation | Predecessor | Complete |

---

## Definition of Done

- All acceptance criteria are met
- Code reviewed and approved
- Committed to `main` via approved PR
