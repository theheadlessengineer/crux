# US-1115 — Generate Cost Allocation, Data Classification, and Compliance Stubs

**Epic:** 1.1 Tier 1 Standards Generation
**Phase:** 1 — Pilot
**Priority:** Must Have
**Status:** To Do

---

## User Story

As a user of Crux,
I want every generated service to include cost allocation tags, a data classification stub, a log retention declaration, and a service catalog entry,
so that FinOps, compliance, and platform teams have structured metadata about every service from its first commit.

---

## Pre-Development Checklist

- [ ] The company's required cost allocation tag keys are agreed (team, service, environment, cost-centre)
- [ ] The data classification levels are agreed (public, internal, confidential, restricted)
- [ ] The service catalog schema is agreed (catalog-entry.yaml format)
- [ ] Epic 1.4 Template Engine is in progress or complete
- [ ] Story estimated and accepted into the sprint

---

## Scope

Generate four YAML stub files as part of every service skeleton. These files are populated with values from the `crux new` prompts and contain documentation comments instructing the team to review and complete them.

### Files Generated

| File | Purpose |
|---|---|
| `cost-budget.yaml` | Expected monthly cloud spend and team attribution |
| `data-classification.yaml` | PII/PHI data fields and classification levels |
| `log-retention.yaml` | Log retention period by environment |
| `catalog-entry.yaml` | Service catalog metadata for platform discoverability |

### In Scope

- Templates for all four files with the service name, team, and environment pre-populated from `crux new` prompts
- Documentation comments in each file explaining what must be completed before production deployment
- Validation that these files exist as part of `crux validate` (Epic 2.5)

### Out of Scope

- Automated cost budget enforcement (Epic 3.9)
- Full GDPR erasure handler (Epic 3.5)
- Log retention enforcement in infrastructure (Epic 3.5)

---

## Acceptance Criteria

- [ ] All four YAML files are present in every generated service
- [ ] Service name and team are pre-populated from `crux new` prompts
- [ ] Each file contains documentation comments identifying what must be completed
- [ ] YAML files are syntactically valid (pass `yamllint`)
- [ ] Files are located under `docs/` or `infra/compliance/` as agreed by the team

---

## Post-Completion Checklist

- [ ] Code reviewed by at least one other platform engineer
- [ ] YAML syntax validated with `yamllint`
- [ ] Content reviewed with FinOps and Compliance stakeholders
- [ ] Story moved to Done in the project tracker

---

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| Epic 1.4 Template Engine | Predecessor | Required |
| Epic 1.3 Prompt Engine | Predecessor | Service name and team captured in prompts |

---

## Definition of Done

- All acceptance criteria are met
- YAML files validated
- Code reviewed and approved
- Committed to `main` via approved PR
