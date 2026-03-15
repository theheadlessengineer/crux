# US-1116 — Generate CI Compliance Checks (SBOM, Licence Scan, Pre-Commit, DAST Slot)

**Epic:** 1.1 Tier 1 Standards Generation
**Phase:** 1 — Pilot
**Priority:** Must Have
**Status:** To Do

---

## User Story

As a user of Crux,
I want every generated service to include CI pipeline steps for SBOM generation, dependency licence scanning, and a disabled DAST slot,
so that supply chain security checks are built into the pipeline from the first commit rather than added as an afterthought.

---

## Pre-Development Checklist

- [ ] The SBOM generation tool is agreed (Syft or equivalent)
- [ ] The licence scanning tool is agreed (FOSSA, Licensee, or equivalent)
- [ ] The DAST tool is agreed (OWASP ZAP or equivalent) — must be disabled by default with a clear comment
- [ ] Epic 1.4 Template Engine is in progress or complete
- [ ] Story estimated and accepted into the sprint

---

## Scope

Generate a GitHub Actions CI pipeline (or the CI platform in use) with the required compliance steps integrated.

### Required CI Steps

| Step | Tool | Default State |
|---|---|---|
| SBOM generation | Syft | Enabled — runs on every build |
| SBOM attestation upload | SLSA / cosign | Enabled |
| Dependency licence scan | Agreed tool | Enabled — fails on GPL-licensed dependencies |
| Pre-commit hook execution | Shell | Enabled |
| DAST scan slot | OWASP ZAP | Disabled — stubbed with `if: false` and a TODO comment |

### In Scope

- CI workflow YAML with all steps
- SBOM artifact uploaded as a build artifact
- Licence scan failing the build on disallowed licences (GPL, AGPL by default)
- DAST step present but explicitly disabled with instructions for enabling it
- Pre-commit hooks run in CI (format, lint, secret scan)

### Out of Scope

- Secret scanning setup (generated in US-1104 pre-commit hooks)
- SAST tooling (not in scope for Phase 1)
- Container scanning (later phase)

---

## Acceptance Criteria

- [ ] CI pipeline includes SBOM generation step that runs on every build
- [ ] SBOM artifact is uploaded and retained for 30 days
- [ ] Licence scan step runs and fails the build if disallowed licences are present
- [ ] DAST step is present but disabled with a `TODO` comment explaining how to enable it
- [ ] Pre-commit hook step runs in CI and fails on violations
- [ ] Pipeline passes on a clean generated service

---

## Post-Completion Checklist

- [ ] Code reviewed by at least one other platform engineer
- [ ] CI pipeline run observed end-to-end on a generated service
- [ ] SBOM artifact confirmed present in build artifacts
- [ ] Deliberately introduced a GPL-licensed dependency — licence scan correctly failed
- [ ] Story moved to Done in the project tracker

---

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| Epic 1.4 Template Engine | Predecessor | Required |
| Epic 0.1 CI pipeline (US-0102) | Predecessor | Base pipeline must exist |

---

## Definition of Done

- All acceptance criteria are met
- Code reviewed and approved
- Committed to `main` via approved PR
