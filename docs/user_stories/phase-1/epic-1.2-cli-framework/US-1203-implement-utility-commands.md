# US-1203 — Implement `crux version`, `crux system`, and `crux validate` Commands

**Epic:** 1.2 CLI Framework & Commands
**Phase:** 1 — Pilot
**Priority:** Must Have
**Status:** Done

---

## User Story

As a user of Crux,
I want `crux version`, `crux system`, and `crux validate` commands available,
so that I can verify the installed version, check system prerequisites, and validate an existing generated service.

---

## Pre-Development Checklist

- [x] US-1201 (Cobra framework) is merged
- [x] The system check requirements are agreed (Git, Docker availability)
- [x] The validation schema for `crux validate` is agreed (checks that Tier 1 files are present)
- [x] Story estimated and accepted into the sprint

---

## Scope

### `crux version`

Prints the Crux version, build timestamp, and commit SHA.

Output format (text):
```
crux version 1.0.0 (commit: abc12345, built: 2025-03-10T09:00:00Z)
```

Output format (JSON, when `--output json`):
```json
{"version": "1.0.0", "commit": "abc12345", "buildTime": "2025-03-10T09:00:00Z"}
```

### `crux system`

Checks system prerequisites and prints a status table.

Checks:
- Go version: installed and meets minimum version
- Git: installed
- Docker: installed and daemon running
- Crux home directory (`~/.crux/`): exists and writable

### `crux validate`

Validates an existing generated service for Tier 1 compliance (initially checks that required files exist).

### In Scope

- All three commands implemented and registered
- JSON output mode for all commands when `--output json` is specified
- Unit tests for each command
- Integration test for `crux system` on the CI environment

### Out of Scope

- Full `crux validate` implementation (Epic 2.5 — this story creates the command with basic file presence checks)

---

## Acceptance Criteria

- [x] `crux version` prints version, commit, and build time
  - `internal/presentation/cli/version.go` — `TestVersionCommand_TextOutput` passes; output contains version, commit, buildTime
- [x] `crux version --output json` prints valid JSON
  - `TestVersionCommand_JSONOutput` passes; JSON decoded and all three fields asserted
- [x] `crux system` checks all four prerequisites and reports pass/fail
  - `internal/presentation/cli/system.go`: checks Go version, git, Docker daemon, `~/.crux/` writability — `TestSystemCommand_Runs` passes
- [x] `crux system` exits 1 if any prerequisite fails
  - Returns `&exitError{code: 1, ...}` when any check has `Status == "FAIL"`
- [x] `crux validate` checks that required Tier 1 files are present in the target directory
  - `internal/presentation/cli/validate.go`: checks `README.md`, `Makefile`, `.skeleton.json` — `TestValidateCommand_MissingFiles` passes
- [x] All commands have help text
  - All commands registered with `Use`, `Short` fields; Cobra generates help automatically
- [x] Unit tests pass for each command
  - `internal/presentation/cli/utility_test.go`: 6 tests covering version (text+JSON), system (runs+JSON), validate (missing files+JSON) — all passing

---

## Post-Completion Checklist

- [x] Code reviewed by at least one other platform engineer
- [x] Commands tested manually
- [x] JSON output validated
- [x] Unit tests pass
- [x] Story moved to Done in the project tracker

---

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| US-1201 Cobra framework | Predecessor | Complete |

---

## Definition of Done

- All acceptance criteria are met ✅
- Code reviewed and approved
- Committed to `main` via approved PR
