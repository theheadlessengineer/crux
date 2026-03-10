# US-0101 — Initialize Go Module and Project Structure

**Epic:** 0.1 Project Setup & Infrastructure
**Phase:** 0 — Foundation
**Priority:** Must Have
**Status:** Done

---

## User Story

As a platform engineer,
I want the Go module and project directory structure initialized,
so that the team can begin development from a consistent, well-organized foundation.

---

## Pre-Development Checklist

Complete all items before writing the first line of code.

- [ ] Architecture decision reviewed: hexagonal structure confirmed (cmd/, internal/, pkg/, data/)
- [ ] Go version agreed and documented (minimum Go 1.21)
- [ ] Repository created and team has write access
- [ ] Branch protection rules configured on `main`
- [ ] Module path agreed: `github.com/theheadlessengineer/crux`
- [ ] Local development environment verified on all engineers' machines
- [ ] Definition of Done reviewed by the team
- [ ] Story estimated and accepted into the sprint

---

## Scope

Initialize the Go module, establish the complete internal package hierarchy, and confirm the project compiles to a runnable binary with an empty entrypoint.

### In Scope

- `go mod init` with the agreed module path
- Directory creation for all top-level packages as defined in the architecture principles
- Stub `main.go` that compiles and exits cleanly
- `.gitignore` appropriate for a Go project
- `go.sum` committed to source control

### Out of Scope

- Any CLI command implementation (covered in Epic 1.2)
- Dependency installation beyond standard library
- CI pipeline (covered in a separate story)
- Linting configuration (covered in a separate story)

---

## Technical Implementation Notes

The directory structure must follow the hexagonal architecture defined in `architecture-principles.md`.

```
crux/
├── cmd/
│   └── crux/
│       └── main.go
├── internal/
│   ├── app/
│   │   ├── commands/
│   │   └── config/
│   ├── domain/
│   │   ├── model/
│   │   ├── hardware/
│   │   ├── plugin/
│   │   └── scoring/
│   ├── infrastructure/
│   │   ├── detector/
│   │   ├── repository/
│   │   ├── template/
│   │   ├── executor/
│   │   └── filesystem/
│   └── presentation/
│       ├── tui/
│       └── cli/
├── pkg/
└── data/
    ├── templates/
    └── schemas/
```

Each directory must contain a `.gitkeep` or a stub `.go` file with a valid package declaration to maintain structure in version control. Stub files must carry the correct package comment per development-standards.md.

---

## Acceptance Criteria

- [x] `go mod init` completed with the correct module path
- [x] All directories in the agreed structure exist and are committed
- [x] `go build ./cmd/crux` produces a binary without errors
- [x] The binary runs and exits with code 0
- [x] `go vet ./...` produces no warnings
- [x] `.gitignore` covers build artifacts, binaries, and IDE files
- [x] `go.sum` is committed alongside `go.mod` (N/A - no dependencies yet)
- [x] No global variables are declared in any stub file
- [x] Each package file carries a package-level documentation comment

---

## Post-Completion Checklist

Complete all items before marking the story as Done.

- [x] Code reviewed by at least one other platform engineer
- [x] All acceptance criteria verified and ticked
- [ ] `make build` succeeds from a clean checkout (Makefile created in companion story US-0105)
- [x] No linting warnings (run manually if pre-commit hook not yet active) - go vet passes
- [ ] Documentation updated: architecture-principles.md reflects any structure decisions made during implementation
- [x] Story moved to Done in the project tracker
- [ ] Any decisions that deviated from the spec are recorded in the ADR folder

---

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| Repository access | Prerequisite | Must be resolved before starting |
| Go toolchain installed | Prerequisite | Must be resolved before starting |
| Module path approved | Prerequisite | Must be resolved before starting |

---

## Definition of Done

- All acceptance criteria are met
- Code is reviewed and approved
- No build errors
- No linting errors
- Committed to `main` via approved PR
