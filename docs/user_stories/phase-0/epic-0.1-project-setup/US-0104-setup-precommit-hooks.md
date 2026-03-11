# US-0104 — Set Up Pre-Commit Hooks

**Epic:** 0.1 Project Setup & Infrastructure
**Phase:** 0 — Foundation
**Priority:** Must Have
**Status:** To Do

---

## User Story

As a platform engineer,
I want pre-commit hooks installed and enforced locally,
so that formatting errors, lint violations, and failing tests are caught before a commit reaches the remote repository.

---

## Pre-Development Checklist

- [ ] US-0103 (Configure linting and formatting) is merged
- [ ] Team agrees on the hook tooling approach (shell script in `.git/hooks/` or a hook manager)
- [ ] All engineers' machines have `gofmt` and `golangci-lint` available in PATH
- [ ] Story estimated and accepted into the sprint

---

## Scope

Implement pre-commit hooks that enforce the three gates: formatting, linting, and test passage. Provide a setup script or Makefile target that installs the hooks into a local clone.

### In Scope

- Pre-commit hook script covering `gofmt`, `golangci-lint run`, and `go test ./...`
- A Makefile target (`make hooks`) or setup script to install the hooks
- Documentation in the README on how to install hooks after cloning
- The hook must exit non-zero to block the commit if any gate fails

### Out of Scope

- Secret scanning hook (may be added in Epic 1.1)
- Commit message linting (covered in conventional commits config)
- Remote-side enforcement (that is the CI pipeline's role)

---

## Technical Implementation Notes

The hook script must follow the pattern in development-standards.md:

```bash
#!/bin/bash
# .git/hooks/pre-commit

if [ -n "$(gofmt -l .)" ]; then
    echo "Code is not formatted. Run: gofmt -w ."
    exit 1
fi

golangci-lint run
if [ $? -ne 0 ]; then
    exit 1
fi

go test ./...
if [ $? -ne 0 ]; then
    exit 1
fi
```

The script must be executable (`chmod +x`). The Makefile target must copy the script into `.git/hooks/pre-commit` and set the executable bit.

Because `.git/hooks/` is not tracked by Git, the canonical script must live in the repository (for example, `scripts/hooks/pre-commit`) and the Makefile copies it. This is the approach that survives re-cloning.

---

## Acceptance Criteria

- [ ] The hook script exists at `scripts/hooks/pre-commit` in the repository
- [ ] `make hooks` (or equivalent) installs the hook into `.git/hooks/pre-commit`
- [ ] A commit with an unformatted file is blocked by the hook
- [ ] A commit with a lint violation is blocked by the hook
- [ ] A commit with a failing test is blocked by the hook
- [ ] A commit with clean code, no lint issues, and passing tests succeeds
- [ ] The README contains instructions for running `make hooks` after cloning
- [ ] The hook script is executable after installation

---

## Post-Completion Checklist

- [ ] Code reviewed by at least one other platform engineer
- [ ] Each engineer on the team has installed the hooks and verified they work
- [ ] All three blocking scenarios tested manually (bad format, bad lint, failing test)
- [ ] README instructions verified by a team member who followed them from scratch
- [ ] Story moved to Done in the project tracker

---

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| US-0103 Configure linting and formatting | Predecessor | Must be merged first |
| US-0105 Create Makefile | Parallel | Coordinate target naming |

---

## Definition of Done

- All acceptance criteria are met
- All engineers have installed and verified the hooks
- Code reviewed and approved
- Committed to `main` via approved PR
