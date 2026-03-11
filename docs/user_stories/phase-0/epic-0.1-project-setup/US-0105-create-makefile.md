# US-0105 — Create Makefile with Common Commands

**Epic:** 0.1 Project Setup & Infrastructure
**Phase:** 0 — Foundation
**Priority:** Must Have
**Status:** To Do

---

## User Story

As a platform engineer,
I want a Makefile with common development commands,
so that every engineer uses consistent commands regardless of their IDE or shell configuration.

---

## Pre-Development Checklist

- [ ] US-0101 (Initialize Go module) is merged
- [ ] Team agrees on the required Makefile targets
- [ ] GNU Make availability confirmed on all developer machines
- [ ] Story estimated and accepted into the sprint

---

## Scope

Create a root-level `Makefile` with the agreed standard targets. The Makefile must be self-documenting via a `help` target.

### In Scope

The following targets are required at minimum for the Foundation phase:

| Target | Description |
|---|---|
| `make build` | Compile the binary to `./bin/crux` |
| `make test` | Run all tests with the race detector |
| `make lint` | Run golangci-lint |
| `make fmt` | Format all Go files |
| `make clean` | Remove build artifacts |
| `make hooks` | Install pre-commit hooks |
| `make help` | Print available targets with descriptions |
| `make vet` | Run go vet |
| `make coverage` | Run tests and open coverage report |

### Out of Scope

- Release and publish targets (later phase)
- `make dev` (requires Docker Compose, Epic 2.9)
- `make upgrade` (Epic 3.1)
- Language-specific targets beyond Go

---

## Technical Implementation Notes

The binary output path must be `./bin/crux`. The `bin/` directory must be listed in `.gitignore`.

The `help` target must parse targets with double-hash comments (`##`) and print them formatted. This is a self-documenting Makefile pattern:

```makefile
.PHONY: help
help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	  awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'
```

The `build` target must pass build flags including version information via `-ldflags`:

```makefile
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

build: ## Compile binary to ./bin/crux
	go build $(LDFLAGS) -o ./bin/crux ./cmd/crux
```

All targets that do not produce files must be declared `.PHONY`.

---

## Acceptance Criteria

- [ ] `make build` produces `./bin/crux` binary
- [ ] `make test` runs all tests with race detection
- [ ] `make lint` runs golangci-lint and returns exit code 0 on clean code
- [ ] `make fmt` formats all Go files in place
- [ ] `make clean` removes `./bin/` and any coverage artifacts
- [ ] `make hooks` installs the pre-commit hook
- [ ] `make help` prints all targets with descriptions
- [ ] `make vet` runs go vet with no warnings
- [ ] All phony targets are declared `.PHONY`
- [ ] Makefile works correctly from a clean repository checkout
- [ ] Version information is embedded in the binary at build time

---

## Post-Completion Checklist

- [ ] Code reviewed by at least one other platform engineer
- [ ] Each target tested individually from a clean checkout
- [ ] `make help` output reviewed for clarity
- [ ] CI pipeline updated to use Makefile targets where appropriate
- [ ] Story moved to Done in the project tracker

---

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| US-0101 Initialize Go module | Predecessor | Must be merged first |
| US-0103 Configure linting | Predecessor | `make lint` depends on this |
| US-0104 Pre-commit hooks | Parallel | Coordinate `make hooks` target |

---

## Definition of Done

- All acceptance criteria are met
- All targets verified working from a clean checkout
- Code reviewed and approved
- Committed to `main` via approved PR
