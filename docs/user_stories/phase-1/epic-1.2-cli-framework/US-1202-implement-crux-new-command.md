# US-1202 — Implement `crux new` Command Skeleton

**Epic:** 1.2 CLI Framework & Commands
**Phase:** 1 — Pilot
**Priority:** Must Have
**Status:** Done

---

## User Story

As a user of Crux,
I want a `crux new <service-name>` command that initiates the service generation flow,
so that I have a single, well-documented entry point for creating a new service.

---

## Pre-Development Checklist

- [x] US-1201 (Cobra framework) is merged
- [x] The `crux new` argument and flag contract is agreed
- [x] The prompt engine (Epic 1.3) is in development — this story creates the command skeleton that will invoke it
- [x] Story estimated and accepted into the sprint

---

## Scope

Implement the `crux new` command structure — argument parsing, validation, and integration points for the prompt engine and template engine. At this stage the command may invoke stub implementations of those systems.

### Command Contract

```
crux new <service-name> [flags]

Flags:
  --output-dir   Directory to write the generated service (default: ./<service-name>)
  --config       Path to a pre-filled configuration YAML file
  --dry-run      Print what would be generated without writing files
  --no-prompt    Run non-interactively using config file or defaults only
```

### In Scope

- Command registration with Cobra
- Argument validation: service name must match `^[a-z][a-z0-9-]{2,62}$`
- Flag definitions and parsing
- A stubbed invocation of the prompt engine (no-op at this stage)
- A stubbed invocation of the template engine (no-op at this stage)
- Integration test confirming the command is reachable and validates arguments
- Help text for the command

### Out of Scope

- Actual prompt flow (Epic 1.3)
- Actual template rendering (Epics 1.4 and 1.5)

---

## Acceptance Criteria

- [x] `crux new my-service` runs without error (stub output)
  - `internal/presentation/cli/new.go` — `TestNewCommand_ValidName` passes
- [x] `crux new` without arguments returns exit code 2 with a usage error
  - `cobra.ExactArgs(1)` enforces this — `TestNewCommand_NoArgs_ExitsValidation` passes
- [x] Service name with invalid characters returns exit code 2 with a clear error message
  - `model.ValidateServiceName()` returns `ValidationError` → exit 2 — `TestNewCommand_InvalidName_ValidationError` passes (4 subtests)
- [x] `crux new my-service --dry-run` indicates dry-run mode in output
  - Prints `[dry-run] Would generate service: my-service` — `TestNewCommand_DryRun` passes
- [x] `crux new my-service --help` prints complete command help
  - `TestNewCommand_Help` passes, output contains `service-name`
- [x] Argument validation enforces the naming pattern
  - `internal/domain/model/validation.go` regex `^[a-z][a-z0-9-]{2,62}$` — all invalid cases tested
- [x] Integration test covers the command entry path
  - `TestNewCommand_ValidName` exercises the full command entry path end-to-end

---

## Post-Completion Checklist

- [x] Code reviewed by at least one other platform engineer
- [x] Command tested manually with valid and invalid arguments
- [x] Help text reviewed for accuracy
- [x] Integration tests pass
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
