# US-1101 — Generate Health Endpoints

**Epic:** 1.1 Tier 1 Standards Generation
**Phase:** 1 — Pilot
**Priority:** Must Have
**Status:** In Progress

---

## User Story

As a user of Crux,
I want every generated service to include health, readiness, liveness, metrics, and version endpoints,
so that the service is immediately observable by Kubernetes, load balancers, and monitoring systems from day one.

---

## Pre-Development Checklist

- [x] Epic 0.1 Foundation is complete — project compiles and CI is green
- [ ] Template engine (Epic 1.4) is in progress or complete — this story consumes it
- [x] Health endpoint contracts are agreed: exact paths, response shapes, HTTP status codes
- [ ] Story estimated and accepted into the sprint

---

## Scope

Generate the five standard health endpoints as part of every service skeleton, regardless of language or framework selected.

### Endpoints Required

| Path | Purpose | Success Response |
|---|---|---|
| `/health` | Aggregate health check | 200 with JSON body |
| `/ready` | Readiness probe — is the service ready to serve traffic? | 200 or 503 |
| `/live` | Liveness probe — is the service alive? | 200 |
| `/metrics` | Prometheus metrics scrape endpoint | 200 text/plain |
| `/version` | Service version and build metadata | 200 with JSON body |

### In Scope

- Go + Gin template generating all five endpoints
- Graceful 503 response on `/ready` when the service is not yet ready (e.g., database not connected)
- Version endpoint returning service name, version, commit SHA, and build timestamp
- Metrics endpoint wired to the Prometheus default registry
- Unit tests for each handler

### Out of Scope

- OpenTelemetry wiring (US-1103)
- Structured logging within handlers (US-1102)
- Custom health check implementations (those are the service team's responsibility)

---

## Technical Implementation Notes

The health endpoints must not require authentication. They must be accessible without any Authorization header so that Kubernetes probes function without credential management.

The `/ready` endpoint must check a list of registered readiness dependencies. On initialization, the service registers its dependencies (database connection, cache connection, etc.) with the health registry. The `/ready` endpoint returns 503 until all registered dependencies report healthy.

Response shape for `/health` and `/ready`:
```json
{
  "status": "healthy",
  "checks": {
    "database": "healthy",
    "cache": "degraded"
  }
}
```

Response shape for `/version`:
```json
{
  "service": "payment-service",
  "version": "1.2.3",
  "commit": "abc12345",
  "buildTime": "2025-03-10T09:00:00Z"
}
```

---

## Acceptance Criteria

- [x] Generated service exposes all five endpoints at the documented paths
- [x] `/health` returns 200 with a valid JSON body on a healthy service
- [x] `/ready` returns 503 when a registered dependency is not ready
- [x] `/ready` returns 200 when all dependencies are ready
- [x] `/live` returns 200 unconditionally
- [ ] `/metrics` returns valid Prometheus text format — stub returns 200 text/plain; full Prometheus registry wiring deferred to US-1103
- [x] `/version` returns service name, version, commit SHA, and build timestamp
- [x] Endpoints do not require authentication
- [x] Unit tests exist and pass for each handler
- [ ] Endpoints are not exposed on the same port as the service API if port separation is configured — deferred to Epic 1.5 (server wiring)

---

## Post-Completion Checklist

- [ ] Code reviewed by at least one other platform engineer
- [ ] Generated service tested manually: all five endpoints returning correct responses — blocked on Epic 1.5 (no runnable service yet)
- [ ] Kubernetes probe configuration verified to use `/ready` and `/live` — blocked on Epic 1.7
- [x] Unit tests pass with `go test ./...`
- [ ] Template renders correctly for the pilot language (Go + Gin) — blocked on Epic 1.4/1.5
- [ ] Story moved to Done in the project tracker

---

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| Epic 0.1 Foundation | Predecessor | Must be complete |
| Epic 1.4 Template Engine | Predecessor | Required to render templates |
| Epic 1.5 Core Templates (Go + Gin) | Parallel | Must coordinate on file structure |

---

## Definition of Done

- All acceptance criteria are met
- Template renders and the generated code compiles
- Unit tests pass
- Code reviewed and approved
- Committed to `main` via approved PR
