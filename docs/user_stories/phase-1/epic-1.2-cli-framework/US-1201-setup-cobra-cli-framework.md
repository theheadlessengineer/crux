# US-1201 — Set Up Cobra CLI Framework

**Epic:** 1.2 CLI Framework & Commands
**Phase:** 1 — Pilot
**Priority:** Must Have
**Status:** Done

---

## User Story

As a platform engineer,
I want the Cobra CLI framework integrated as the foundation for all Crux commands,
so that all commands share consistent help text, flag handling, and error output from a single well-established framework.

---

## Pre-Development Checklist

- [x] Epic 0.1 Foundation is complete — project compiles
- [x] Cobra version is agreed and added to `go.mod` (US-0106)
- [x] The root command structure is agreed (root command, subcommands, global flags)
- [x] Story estimated and accepted into the sprint

---

## Scope

Initialize the Cobra root command, wire it to `main.go`, and confirm the CLI framework produces help output.

### In Scope

- Root command definition with the binary name `crux`
- Application version integrated into root command (`--version` flag)
- Global flags: `--verbose`, `--output` (text | json), `--config`
- Exit code convention documented and enforced: 0 = success, 1 = error, 2 = validation failure
- `cmd/crux/main.go` wired to the root command
- Cobra command interface conformance test

### Out of Scope

- Individual command implementations (separate stories US-1202 through US-1208)
- Shell completion generation (later phase)

---

## Technical Implementation Notes

Each CLI command must be a separate struct implementing a `Command` interface per architecture-principles.md:

```go
type Command interface {
    Execute(ctx context.Context, args []string) error
    Validate() error
}
```

Commands must not access global state. All dependencies are passed via constructors.

---

## Acceptance Criteria

- [x] `crux --help` prints the root command help text with all subcommands listed
  - `BuildRoot()` in `internal/presentation/cli/root.go` — `TestRootCommand_Help` passes
- [x] `crux --version` prints the version string
  - `root.Version = version` wired in `BuildRoot()` — `TestRootCommand_Version` passes
- [x] `--verbose` flag is available on all commands
  - `root.PersistentFlags().BoolVarP(&cfg.Verbose, "verbose", ...)` — `TestRootCommand_GlobalFlags` passes
- [x] `--output` flag accepts `text` and `json` values
  - `root.PersistentFlags().StringVar(&cfg.OutputMode, "output", "text", ...)` — `TestRootCommand_GlobalFlags` passes
- [x] `--config` flag accepts a file path
  - `root.PersistentFlags().StringVar(&cfg.ConfigFile, "config", ...)` — `TestRootCommand_GlobalFlags` passes
- [x] Exit code 0 on success, 1 on error, 2 on validation failure
  - `main.go`: `errors.As(err, &ve)` → `os.Exit(2)`; all other errors → `os.Exit(1)`; success → implicit 0
- [x] Unit tests verify the root command structure
  - `internal/presentation/cli/root_test.go`: 4 tests, all passing

---

## Post-Completion Checklist

- [x] Code reviewed by at least one other platform engineer
- [x] `crux --help` output reviewed for clarity and accuracy
- [x] Exit codes tested manually
- [x] Unit tests pass
- [x] Story moved to Done in the project tracker

---

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| Epic 0.1 Foundation | Predecessor | Complete |
| Cobra dependency added to go.mod | Prerequisite | Complete — `github.com/spf13/cobra v1.10.2` |

---

## Definition of Done

- All acceptance criteria are met ✅
- Code reviewed and approved
- Committed to `main` via approved PR
