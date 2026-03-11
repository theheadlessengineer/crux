# US-0106 — Configure Dependency Management (go.mod)

**Epic:** 0.1 Project Setup & Infrastructure
**Phase:** 0 — Foundation
**Priority:** Must Have
**Status:** To Do

---

## User Story

As a platform engineer,
I want the Go module's dependency management configured and documented,
so that all engineers use consistent, pinned dependencies and the team has a clear process for adding or updating them.

---

## Pre-Development Checklist

- [ ] US-0101 (Initialize Go module) is merged
- [ ] Team agrees on the process for approving new dependencies
- [ ] Licences for any immediately required dependencies are reviewed
- [ ] Story estimated and accepted into the sprint

---

## Scope

Establish the baseline `go.mod` and `go.sum` files, add the initial set of approved direct dependencies, and document the dependency management process.

### Initial Required Dependencies

| Package | Purpose |
|---|---|
| `github.com/spf13/cobra` | CLI framework (Epic 1.2) — add now to unblock parallel work |
| `github.com/stretchr/testify` | Test assertions |

No other dependencies should be added at Foundation phase without team review.

### In Scope

- `go.mod` with the agreed initial dependency set
- `go.sum` committed and valid
- Documentation in the contribution guide on the process to add new dependencies
- Dependency licence review recorded for the initial set

### Out of Scope

- Adding all Epic 1 dependencies (those are added by the respective stories)
- Automated licence scanning CI step (Epic 1.1)

---

## Technical Implementation Notes

Running `go mod tidy` must be part of the development workflow and the CI pipeline. It must be run after any dependency change.

The `go.sum` file must always be committed. Engineers must not add the `vendor/` directory unless the team explicitly decides to vendor dependencies, and that decision must be recorded in an ADR.

Dependency updates follow the process:
1. Engineer proposes update with reasoning in the PR description
2. Licence of the new or updated package is verified
3. At least one reviewer confirms the licence is acceptable
4. `go mod tidy` is run before the PR is merged

---

## Acceptance Criteria

- [ ] `go.mod` contains the correct module path and Go version
- [ ] Initial dependencies (`cobra`, `testify`) are present and pinned in `go.sum`
- [ ] `go mod tidy` produces no changes after a clean checkout
- [ ] `go.sum` is committed and consistent with `go.mod`
- [ ] Dependency management process documented in `CONTRIBUTING.md`
- [ ] Licences for initial dependencies are recorded

---

## Post-Completion Checklist

- [ ] Code reviewed by at least one other platform engineer
- [ ] `go mod tidy` verified to produce no changes on a clean checkout
- [ ] Contributing guide reviewed for clarity
- [ ] Story moved to Done in the project tracker

---

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| US-0101 Initialize Go module | Predecessor | Must be merged first |

---

## Definition of Done

- All acceptance criteria are met
- Code reviewed and approved
- Committed to `main` via approved PR
