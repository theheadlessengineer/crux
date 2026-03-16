# US-1601 — Define Plugin Manifest Schema and Implement Plugin Loader

**Epic:** 1.6 Plugin System — Core Infrastructure
**Phase:** 1 — Pilot
**Priority:** Must Have
**Status:** Done

---

## User Story

As a platform engineer building Crux,
I want a plugin manifest schema and loader that validates and loads plugins at startup,
so that the core system can discover and initialize plugins in a consistent, safe way.

---

## Pre-Development Checklist

- [x] The plugin manifest schema is agreed (`plugin.yaml` fields and types)
- [x] The plugin trust tier model is agreed (Tier 1 = official, Tier 2 = inner source, Tier 3 = community)
- [x] The plugin isolation strategy is agreed (in-process with interface constraints at MVP)
- [ ] Epic 1.4 Template Engine is merged
- [x] Story estimated and accepted into the sprint

---

## Scope

Define the `plugin.yaml` manifest schema, implement the manifest parser and validator, and implement the plugin loader that discovers and initializes plugins from the filesystem.

### plugin.yaml Schema

```yaml
apiVersion: crux/v1
kind: Plugin
metadata:
  name: crux-plugin-postgresql
  version: 1.0.0
  description: PostgreSQL integration with connection pooling and migrations
  author: Platform Engineering
  trustTier: 1
  tags: [database, postgresql]
  cruxVersionConstraint: ">=1.0.0"
spec:
  questions:
    - id: pg_version
      type: select
      prompt: "Which PostgreSQL version?"
      options: ["15", "16"]
      default: "16"
  templates: []
  hooks:
    preGenerate: []
    postGenerate: []
  dependencies: []
```

### In Scope

- `plugin.yaml` schema defined as Go structs with YAML tags
- Manifest parser loading from a file path
- Manifest validator checking required fields, version format, and trust tier
- Plugin loader that discovers all plugins from `~/.crux/plugins/` and the embedded official plugins directory
- Plugin version compatibility check against the running Crux version
- Plugin lifecycle hooks (`preGenerate`, `postGenerate`) invoked at the correct points
- Unit tests for parsing, validation, loading, and hook execution

### Out of Scope

- Plugin registry download (Epic 2.4)
- Plugin sandbox/process isolation (MVP uses in-process — documented as a known limitation in an ADR)

---

## Acceptance Criteria

- [x] Valid `plugin.yaml` files are parsed without error
- [x] Invalid `plugin.yaml` files (missing required fields, bad version format) are rejected with clear errors
- [x] Plugins incompatible with the running Crux version are rejected at startup
- [x] `preGenerate` hooks execute before template rendering
- [x] `postGenerate` hooks execute after template rendering
- [x] Unit tests cover valid plugin, invalid plugin, incompatible version, and hook execution
- [x] Plugin discovery searches both embedded directory and `~/.crux/plugins/`

---

## Post-Completion Checklist

- [ ] Code reviewed by at least one other platform engineer
- [ ] A test plugin loaded successfully and hooks verified
- [ ] An invalid plugin correctly rejected
- [x] Unit tests pass
- [x] Plugin isolation strategy documented in an ADR
- [ ] Story moved to Done in the project tracker

---

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| Epic 1.4 Template Engine | Predecessor | Required |
| Epic 1.3 Prompt Engine | Predecessor | Plugin questions integrate with prompt engine |

---

## Definition of Done

- All acceptance criteria are met
- ADR documenting plugin isolation strategy committed
- Code reviewed and approved
- Committed to `main` via approved PR
