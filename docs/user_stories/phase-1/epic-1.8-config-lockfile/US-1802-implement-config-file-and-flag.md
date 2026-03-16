# US-1802 — Implement Configuration File Input and --config Flag

**Epic:** 1.8 Configuration & Lockfile System
**Phase:** 1 — Pilot
**Priority:** Must Have
**Status:** To Do

---

## User Story

As a user of Crux,
I want to pass a pre-filled YAML configuration file to `crux new`,
so that I can generate services non-interactively or pre-fill answers for the interactive flow.

---

## Pre-Development Checklist

- [ ] US-1801 (skeleton.json and lockfile) is merged
- [ ] The configuration file schema (YAML) is agreed and documented
- [ ] The precedence order is agreed: CLI flags > config file > defaults
- [ ] Story estimated and accepted into the sprint

---

## Scope

Implement a YAML configuration file format that pre-fills prompt answers, and wire the `--config` flag on `crux new` to load it.

### Configuration File Schema

```yaml
service:
  name: payment-service
  language: go
  framework: gin
  team: payments-team
  environment: production

plugins:
  - crux-plugin-postgresql
  - crux-plugin-redis

answers:
  pg_version: "16"
  redis_mode: cluster
```

### In Scope

- YAML configuration file parser with schema validation
- `--config` flag on `crux new` wired to the parser
- Pre-filled answers skip the corresponding prompts in the interactive flow
- `--no-prompt` flag runs entirely from config file without interactive prompts (all required answers must be present)
- CLI flag precedence: CLI flags override config file values
- Configuration validation fails fast with clear errors for missing required fields when `--no-prompt` is used
- Unit tests for parsing, validation, and precedence

---

## Acceptance Criteria

- [ ] `crux new my-service --config myconfig.yaml` loads the configuration file
- [ ] Pre-filled answers in the config file skip the corresponding prompts
- [ ] `--no-prompt` combined with a complete config file generates a service without any interactive prompts
- [ ] `--no-prompt` without a config file fails with a clear error listing missing required answers
- [ ] CLI flags override config file values when both are present
- [ ] Invalid configuration file (bad YAML, missing required fields in `--no-prompt` mode) fails with a clear error
- [ ] Unit tests verify parsing, validation, and precedence rules

---

## Post-Completion Checklist

- [ ] Code reviewed by at least one other platform engineer
- [ ] Config file tested manually: interactive mode with pre-fills, and `--no-prompt` mode
- [ ] Precedence verified: CLI flag overrides config file value
- [ ] Unit tests pass
- [ ] Story moved to Done in the project tracker

---

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| US-1801 skeleton.json and lockfile | Predecessor | Must be merged |
| US-1202 `crux new` command | Predecessor | `--config` flag must be wired |

---

## Definition of Done

- All acceptance criteria are met
- Code reviewed and approved
- Committed to `main` via approved PR
