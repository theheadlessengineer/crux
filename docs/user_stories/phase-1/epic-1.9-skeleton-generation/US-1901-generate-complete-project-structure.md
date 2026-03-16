# US-1901 — Generate Complete Service Directory Structure

**Epic:** 1.9 Complete Skeleton Generation
**Phase:** 1 — Pilot
**Priority:** Must Have
**Status:** To Do

---

## User Story

As a user of Crux,
I want `crux new` to generate the complete directory structure for a production-ready service,
so that the entire project scaffold — from application code to infrastructure to documentation — is present from the first commit.

---

## Pre-Development Checklist

- [ ] Epics 1.1 through 1.8 are merged or nearing completion
- [ ] The complete directory structure is agreed and documented
- [ ] All template and plugin stories are merged
- [ ] Story estimated and accepted into the sprint

---

## Scope

This story is the integration point: wire all templates, plugins, and Tier 1 standard generators into the `crux new` flow to produce a complete, end-to-end service skeleton.

### Generated Directory Structure

```
<service-name>/
├── cmd/<service-name>/main.go
├── internal/
│   ├── app/
│   ├── domain/
│   └── infrastructure/
├── infra/
│   ├── kubernetes/
│   ├── terraform/
│   └── monitoring/
├── tests/
│   ├── unit/
│   └── integration/
├── scripts/
│   ├── seed.sh
│   ├── check_env.sh
│   ├── snapshot-db.sh
│   └── restore-db.sh
├── docs/
│   ├── adr/
│   ├── runbook.md
│   └── capacity-model.md
├── .skeleton.json
├── crux.lock
├── Makefile
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── .gitignore
├── .editorconfig
├── .envrc
└── README.md
```

### In Scope

- Complete directory structure generation
- All Tier 1 standard files present
- All selected plugin files present
- `docker-compose.yml` for local development
- Scripts directory with seed, environment check, and database management scripts
- Integration test: generated skeleton compiles, runs, and all generated tests pass

---

## Acceptance Criteria

- [ ] Complete directory structure is generated as documented
- [ ] Generated service compiles with `go build ./...`
- [ ] Generated service starts with `make dev`
- [ ] All generated tests pass with `go test ./...`
- [ ] Docker build succeeds
- [ ] All Tier 1 standard files are present
- [ ] `.skeleton.json` and `crux.lock` are generated
- [ ] Documentation stubs are present and non-empty
- [ ] Scripts are executable

---

## Post-Completion Checklist

- [ ] Code reviewed by at least one other platform engineer
- [ ] Generated service built, run, and all endpoints verified
- [ ] All generated tests passed
- [ ] Complete directory structure verified against the agreed specification
- [ ] Story moved to Done in the project tracker

---

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| Epics 1.1 — 1.8 | Predecessors | All must be merged or stable |

---

## Definition of Done

- All acceptance criteria are met
- Complete service generated, compiled, and running
- Code reviewed and approved
- Committed to `main` via approved PR
