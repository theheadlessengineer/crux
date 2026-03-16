# US-1801 — Implement .skeleton.json and crux.lock Generation

**Epic:** 1.8 Configuration & Lockfile System
**Phase:** 1 — Pilot
**Priority:** Must Have
**Status:** To Do

---

## User Story

As a user of Crux,
I want every generated service to include a `.skeleton.json` and `crux.lock` file,
so that the exact decisions and plugin versions used to generate the service are recorded and can be used for future upgrades and audits.

---

## Pre-Development Checklist

- [ ] The `.skeleton.json` schema is agreed (all prompt answers, plugin versions, deviations, metadata)
- [ ] The `crux.lock` schema is agreed (plugin name and exact version only)
- [ ] Epic 1.3 Prompt Engine is merged — answers must be available
- [ ] Epic 1.6 Plugin System is merged — plugin versions must be available
- [ ] Story estimated and accepted into the sprint

---

## Scope

Implement the generation of `.skeleton.json` and `crux.lock` at the end of every `crux new` run.

### .skeleton.json Schema

```json
{
  "cruxVersion": "1.0.0",
  "generatedAt": "2025-03-10T09:00:00Z",
  "service": {
    "name": "payment-service",
    "language": "go",
    "framework": "gin"
  },
  "answers": {},
  "plugins": [
    {"name": "crux-plugin-postgresql", "version": "1.2.0"}
  ],
  "deviations": [],
  "tier1Standards": {
    "enforced": true,
    "disabledStandards": []
  }
}
```

### crux.lock Schema

```json
{
  "lockfileVersion": 1,
  "plugins": {
    "crux-plugin-postgresql": "1.2.0",
    "crux-plugin-redis": "0.9.1"
  }
}
```

### In Scope

- Generation of both files at the end of `crux new`
- All prompt answers recorded in `.skeleton.json`
- All plugin versions recorded in both files
- Deviation field ready for use by `crux audit` and the deviation workflow
- Unit tests verifying schema correctness and generation

### Out of Scope

- `crux upgrade` reading the lockfile (Epic 3.1)
- `crux audit` reading `.skeleton.json` (Epic 2.6)
- Configuration file input to `crux new` (US-1802)

---

## Acceptance Criteria

- [ ] `.skeleton.json` is generated in the service root after every `crux new` run
- [ ] `crux.lock` is generated in the service root after every `crux new` run
- [ ] All prompt answers are present in `.skeleton.json`
- [ ] All plugin versions are present in both files
- [ ] Both files are valid JSON (validated with `jq .`)
- [ ] `generatedAt` timestamp is correct
- [ ] `cruxVersion` reflects the actual Crux version used
- [ ] Unit tests verify schema correctness

---

## Post-Completion Checklist

- [ ] Code reviewed by at least one other platform engineer
- [ ] Both files generated and verified on a real `crux new` run
- [ ] JSON validated with `jq`
- [ ] Unit tests pass
- [ ] Story moved to Done in the project tracker

---

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| Epic 1.3 Prompt Engine | Predecessor | Answers must be available |
| Epic 1.6 Plugin System | Predecessor | Plugin versions must be available |

---

## Definition of Done

- All acceptance criteria are met
- Code reviewed and approved
- Committed to `main` via approved PR
