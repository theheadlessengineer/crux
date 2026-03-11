# US-0108 — Create Initial README and Contribution Guidelines

**Epic:** 0.1 Project Setup & Infrastructure
**Phase:** 0 — Foundation
**Priority:** Must Have
**Status:** Done

---

## User Story

As a platform engineer or contributor,
I want a clear README and contribution guide committed to the repository,
so that any engineer who clones the project can understand its purpose and begin contributing without requiring tribal knowledge.

---

## Pre-Development Checklist

- [x] All other Epic 0.1 stories are merged or in progress — README must reflect actual project state
- [x] Team agrees on the content sections required in the README
- [x] Crux brand explanation and positioning reviewed (crux-brand-explanation.md)
- [x] Story estimated and accepted into the sprint

---

## Scope

Create a `README.md` at the repository root and a `CONTRIBUTING.md` that covers the practical workflow for contributors.

### README.md — Required Sections

- Project name and one-sentence description
- What Crux does and why it exists (drawing from crux-readme.md)
- Prerequisites (Go version, required tools)
- Getting started: clone, install hooks, build, run
- Available Makefile targets (reference to `make help`)
- Project structure overview
- Link to CONTRIBUTING.md
- Link to architecture-principles.md
- Maintainers and support contact

### CONTRIBUTING.md — Required Sections

- Code of conduct reference
- Branch naming conventions (from development-standards.md)
- Commit message format (Conventional Commits)
- Pull request process (from development-standards.md)
- Code review guidelines
- Testing requirements
- Dependency management process
- How to run the full local quality gate
- How to install pre-commit hooks

### Out of Scope

- Plugin development guide (Epic 1.6)
- Full user documentation (Epic 2.7)
- Video tutorials (Epic 2.7)

---

## Technical Implementation Notes

The README is the first document any engineer reads. It must be accurate, direct, and free of placeholder content. Sections that are not yet applicable must be clearly marked as forthcoming rather than left blank or filled with stubs.

The README must reflect the current state of the codebase. If the binary does not yet have a meaningful command, the Getting Started section must say so rather than documenting commands that do not exist yet.

Language must be business professional and technically precise. No marketing language in the CONTRIBUTING.md. The README may carry the positioning language from `crux-readme.md` in the introductory section.

---

## Acceptance Criteria

- [x] `README.md` exists at the repository root with all required sections
- [x] `CONTRIBUTING.md` exists at the repository root with all required sections
- [x] All commands documented in the README are verified to work
- [x] No placeholder or Lorem Ipsum content remains
- [x] Links between README.md, CONTRIBUTING.md, and architecture-principles.md are valid
- [x] A new engineer can follow the README and produce a working binary without additional guidance
- [x] Makefile `make help` output is consistent with what is documented in the README

---

## Post-Completion Checklist

- [x] README reviewed by a team member who was not involved in writing it
- [x] That team member confirmed they could follow it from a clean checkout
- [x] Code reviewed by at least one other platform engineer
- [x] All acceptance criteria verified
- [x] Story moved to Done in the project tracker

---

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| US-0101 through US-0107 | Predecessors | ✅ All Done |

---

## Definition of Done

- All acceptance criteria are met
- A team member not involved in writing the documents has verified them end-to-end
- Code reviewed and approved
- Committed to `main` via approved PR
