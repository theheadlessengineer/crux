# US-1704 — Create `crux-plugin-auth-jwt` Plugin

**Epic:** 1.7 Pilot Plugins (9 Essential Plugins)
**Phase:** 1 — Pilot
**Priority:** Must Have
**Status:** Done

---

## User Story

As a user of Crux,
I want the `crux-plugin-auth-jwt` plugin available and installable,
so that I can generate a service with Auth JWT (JWT Validation, RBAC/ABAC Stubs) integration pre-configured and production-ready.

---

## Pre-Development Checklist

- [x] Epic 1.6 Plugin System (US-1601) is merged — plugin manifest and loader must be in place
- [x] Epic 1.5 Core Templates are merged — the plugin adds to an existing base skeleton
- [x] The plugin's questions, generated files, and integration points are documented
- [x] Integration-specific defaults and configuration options are agreed
- [x] Story estimated and accepted into the sprint

---

## Scope

Implement the `crux-plugin-auth-jwt` plugin including its `plugin.yaml` manifest, prompt questions, templates, and lifecycle hooks.

### Plugin Deliverables

- `plugin.yaml` manifest with metadata, version, trust tier, questions, and template references
- All templates required to integrate Auth JWT (JWT Validation, RBAC/ABAC Stubs) into the generated service
- Pre-generate and post-generate hooks if required
- README for the plugin documenting its questions, generated files, and usage
- Unit tests for all templates and hooks
- Integration test confirming the generated service compiles with the plugin applied

### Integration Standards

All generated integration code must:
- Follow the architecture principles (no global state, dependency injection, fail fast on misconfiguration)
- Emit structured log messages at connection, error, and disconnection events
- Register readiness dependencies with the health check registry
- Handle connection errors with configurable retry logic
- Include shutdown integration — connections closed cleanly on SIGTERM

---

## Acceptance Criteria

- [x] `plugin.yaml` manifest is valid and loads without error
- [x] Plugin questions appear correctly in the `crux new` flow
- [x] All templates render without error
- [ ] Generated service with this plugin compiles with `go build ./...`
- [ ] Generated service starts with the integration connected
- [ ] Health endpoint returns unhealthy when the integration is unavailable
- [ ] Shutdown closes the integration connection cleanly
- [x] Unit tests pass
- [ ] Integration test confirms compilation and startup
- [x] Plugin README documents all generated files and configuration options

---

## Post-Completion Checklist

- [ ] Code reviewed by at least one other platform engineer
- [ ] Plugin tested end-to-end: `crux new` flow, compile, run, health check
- [ ] Integration tested against a real Auth JWT (JWT Validation, RBAC/ABAC Stubs) instance in the local environment
- [ ] Shutdown behaviour verified
- [x] Unit tests pass
- [ ] Plugin README reviewed for accuracy
- [ ] Story moved to Done in the project tracker

---

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| Epic 1.6 Plugin System (US-1601) | Predecessor | Must be merged |
| Epic 1.5 Core Templates | Predecessor | Base service must exist |

---

## Definition of Done

- All acceptance criteria are met
- Plugin tested end-to-end
- Code reviewed and approved
- Committed to `main` via approved PR
