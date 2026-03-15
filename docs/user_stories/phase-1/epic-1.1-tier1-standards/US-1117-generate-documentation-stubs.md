# US-1117 — Generate Documentation Stubs (ADR, Runbook, Capacity Model, Conventional Commits)

**Epic:** 1.1 Tier 1 Standards Generation
**Phase:** 1 — Pilot
**Priority:** Must Have
**Status:** To Do

---

## User Story

As a user of Crux,
I want every generated service to include an ADR folder, a runbook template, a capacity model template, and conventional commits configuration,
so that the documentation and decision-recording scaffolding is ready on day one.

---

## Pre-Development Checklist

- [ ] The ADR template format is agreed (Markdown, MADR format or company variant)
- [ ] The runbook template sections are agreed (incident response, disaster recovery section required)
- [ ] The capacity model template is reviewed by the FinOps team
- [ ] Epic 1.4 Template Engine is in progress or complete
- [ ] Story estimated and accepted into the sprint

---

## Scope

Generate the documentation directory structure and seed files as part of every service skeleton.

### Generated Files

| File / Directory | Description |
|---|---|
| `docs/adr/` | Architecture Decision Record folder |
| `docs/adr/ADR-001-initial-technology-choices.md` | First ADR documenting the choices made during `crux new` |
| `docs/runbook.md` | Runbook template with incident and disaster recovery sections |
| `docs/capacity-model.md` | Capacity model template for resource planning |
| `docs/TODO.md` | Intentional placeholder list — things the team must complete before go-live |
| `.commitlintrc.yaml` | Conventional Commits configuration |
| `CHANGELOG.md` | Empty CHANGELOG ready for first entry |
| `.editorconfig` | Editor configuration for consistent formatting |

### In Scope

- All files above with appropriate placeholder content and documentation comments
- ADR-001 auto-populated with the technology decisions made during `crux new` (language, framework, plugins selected)
- TODO.md listing all placeholders that must be resolved before production

### Out of Scope

- Automated CHANGELOG generation (later tooling)
- API documentation (OpenAPI stub is a separate story)

---

## Acceptance Criteria

- [ ] All listed files are present in every generated service
- [ ] ADR-001 is auto-populated with the choices made during `crux new`
- [ ] Runbook template includes an incident response section and a disaster recovery section
- [ ] TODO.md lists all placeholders with instructions for resolving each
- [ ] `.editorconfig` enforces UTF-8, LF line endings, and 4-space indentation for Go
- [ ] `.commitlintrc.yaml` enforces Conventional Commits format
- [ ] Files contain no placeholder Lorem Ipsum text

---

## Post-Completion Checklist

- [ ] Code reviewed by at least one other platform engineer
- [ ] ADR-001 content reviewed — choices from `crux new` correctly reflected
- [ ] Runbook template reviewed by an on-call engineer for completeness
- [ ] Story moved to Done in the project tracker

---

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| Epic 1.4 Template Engine | Predecessor | Required |
| Epic 1.3 Prompt Engine | Predecessor | Choices from `crux new` populate ADR-001 |

---

## Definition of Done

- All acceptance criteria are met
- Code reviewed and approved
- Committed to `main` via approved PR
