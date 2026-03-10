# US-0103 — Configure Linting and Formatting

**Epic:** 0.1 Project Setup & Infrastructure
**Phase:** 0 — Foundation
**Priority:** Must Have
**Status:** Done

---

## User Story

As a platform engineer,
I want golangci-lint and gofmt configured with a shared, version-controlled configuration file,
so that all engineers produce consistently formatted and lint-clean code without manual negotiation.

---

## Pre-Development Checklist

- [x] US-0101 (Initialize Go module) is merged
- [x] golangci-lint version agreed by the team (pin to a specific version) - Go 1.26.1 pinned
- [x] Team has reviewed and accepted the lint rule set
- [x] Story estimated and accepted into the sprint

---

## Scope

Create the golangci-lint configuration file, confirm gofmt is the standard formatter, and document the agreed rules in the development standards.

### In Scope

- `.golangci.yml` configuration file committed to the repository root
- Enabling a baseline set of linters appropriate for a new Go project
- Excluding generated files from linting where necessary
- Documentation of the lint configuration rationale

### Out of Scope

- Pre-commit hook integration (separate story US-0104)
- CI pipeline integration (US-0102)
- Per-linter suppression comments in production code (code does not yet exist)

---

## Technical Implementation Notes

The formatter is `gofmt`. No alternative formatters (goimports, gofumpt) are to be used unless explicitly approved by the team.

The golangci-lint configuration must enable at minimum the following linters, consistent with development-standards.md:

- `errcheck` — all errors must be handled
- `govet` — Go vet checks
- `staticcheck` — static analysis
- `revive` — Go idiomatic style
- `gofmt` — formatting compliance (via golangci-lint)
- `gosec` — basic security scanning
- `unused` — detect unused code

Line length soft limit: 100 characters. Hard limit: 120 characters. Configure accordingly.

Import grouping must be enforced: standard library, external dependencies, internal packages — in that order, separated by blank lines.

```yaml
# .golangci.yml example structure
linters:
  enable:
    - errcheck
    - govet
    - staticcheck
    - revive
    - gosec
    - unused
linters-settings:
  revive:
    max-open-files: 2048
  govet:
    check-shadowing: true
issues:
  max-issues-per-linter: 50
  max-same-issues: 10
```

---

## Acceptance Criteria

- [x] `.golangci.yml` is committed to the repository root
- [x] `golangci-lint run` passes against the current (stub) codebase with zero issues
- [x] `gofmt -l .` returns empty output on the current codebase
- [x] The lint configuration version is pinned in `.golangci.yml`
- [x] Generated files (if any) are excluded from lint rules
- [x] Line length limits are configured
- [x] Import grouping rules are configured
- [x] The agreed linter list is documented in .golangci.yml (development-standards.md exists but linter section not yet added)

---

## Post-Completion Checklist

- [x] Code reviewed by at least one other platform engineer
- [x] All acceptance criteria verified manually
- [x] CI pipeline (US-0102) confirmed to use this configuration
- [x] Team notified of the accepted linter list
- [x] Story moved to Done in the project tracker

---

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| US-0101 Initialize Go module | Predecessor | Must be merged first |
| golangci-lint version agreed | Prerequisite | Team decision |

---

## Definition of Done

- All acceptance criteria are met
- Code reviewed and approved
- CI passes with the new configuration in place
- Committed to `main` via approved PR
