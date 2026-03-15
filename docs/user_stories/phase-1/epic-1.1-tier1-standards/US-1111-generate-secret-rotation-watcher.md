# US-1111 — Generate Secret Rotation Watcher

**Epic:** 1.1 Tier 1 Standards Generation
**Phase:** 1 — Pilot
**Priority:** Must Have
**Status:** To Do

---

## User Story

As a user of Crux,
I want every generated service to include a secret rotation watcher that reloads secrets when they are rotated,
so that services continue operating after a credential rotation without requiring a pod restart.

---

## Pre-Development Checklist

- [ ] US-1110 (Secrets management config) is merged
- [ ] The secrets backend's rotation notification mechanism is understood (Vault lease renewal, AWS Secrets Manager rotation event)
- [ ] The team agrees on which secrets require rotation support at MVP (typically database passwords, API keys)
- [ ] Story estimated and accepted into the sprint

---

## Scope

Generate a background watcher that periodically refreshes secrets from the backend and signals the application to reload affected components.

### In Scope

- A background goroutine that polls the secrets backend on a configurable interval (`SECRETS_REFRESH_INTERVAL_SECONDS`, default: 300)
- Signal propagation to registered components (e.g., database connection pool) when a secret changes
- Structured log message emitted when a secret is rotated
- Graceful shutdown integration — watcher stops cleanly on SIGTERM
- Unit tests using a mock secrets client that simulates rotation

### Out of Scope

- Zero-downtime connection pool drain during rotation (that complexity is handled in the database plugin)
- Push-based rotation triggers (polling is sufficient for MVP)

---

## Acceptance Criteria

- [ ] Background watcher polls secrets on the configured interval
- [ ] When a secret value changes, registered components are notified
- [ ] A structured log message is emitted when rotation is detected
- [ ] Watcher shuts down cleanly on SIGTERM
- [ ] `SECRETS_REFRESH_INTERVAL_SECONDS` controls the polling interval
- [ ] No secret values appear in log messages
- [ ] Unit tests simulate a rotation and verify the notification signal

---

## Post-Completion Checklist

- [ ] Code reviewed by at least one other platform engineer
- [ ] Rotation simulated in test environment — components notified correctly
- [ ] Log output verified to show rotation detection without exposing values
- [ ] Unit tests pass
- [ ] Story moved to Done in the project tracker

---

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| US-1110 Secrets management config | Predecessor | Required |
| US-1105 Graceful shutdown | Predecessor | Watcher must integrate with shutdown |

---

## Definition of Done

- All acceptance criteria are met
- Code reviewed and approved
- Committed to `main` via approved PR
