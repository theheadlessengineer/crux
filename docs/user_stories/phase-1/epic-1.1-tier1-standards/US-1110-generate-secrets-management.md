# US-1110 — Generate Secrets Management Configuration (Vault / AWS Secrets Manager)

**Epic:** 1.1 Tier 1 Standards Generation
**Phase:** 1 — Pilot
**Priority:** Must Have
**Status:** To Do

---

## User Story

As a user of Crux,
I want every generated service to retrieve secrets from Vault or AWS Secrets Manager at startup rather than from environment variables or mounted files,
so that secrets are never stored in container images, Kubernetes ConfigMaps, or version control.

---

## Pre-Development Checklist

- [ ] The company's preferred secrets backend (Vault or AWS Secrets Manager) is agreed and documented
- [ ] Vault or AWS Secrets Manager is operational in the target environments
- [ ] The secret naming convention is agreed (e.g., `<env>/<service-name>/<key>`)
- [ ] Epic 1.4 Template Engine is in progress or complete
- [ ] Story estimated and accepted into the sprint

---

## Scope

Generate a secrets configuration loader that fetches secrets from the agreed backend at service startup and makes them available to the application through a typed configuration object.

### In Scope

- Configuration loader that fetches secrets from Vault (AppRole or Kubernetes auth) or AWS Secrets Manager
- Secrets fetched at startup before the HTTP server starts
- Failure to fetch required secrets results in a startup failure with a clear error message (Fail Fast principle)
- Environment variable `SECRETS_BACKEND` selects the backend (`vault` or `aws-secrets-manager`)
- Generated secrets path stub in the service configuration
- Unit tests using a mock secrets client

### Out of Scope

- Secrets rotation watcher (US-1111 — separate story)
- Encryption key management
- Database password rotation (covered by database plugin stories)

---

## Acceptance Criteria

- [ ] Generated service fetches secrets from the configured backend at startup
- [ ] Service fails to start (exit 1) if required secrets are unavailable
- [ ] Error message clearly identifies which secret is missing
- [ ] No secrets are logged at any log level
- [ ] `SECRETS_BACKEND` environment variable selects Vault or AWS Secrets Manager
- [ ] Unit tests use a mock secrets client and pass
- [ ] Service configuration object is populated from fetched secrets (not from environment variables for secret values)

---

## Post-Completion Checklist

- [ ] Code reviewed by at least one other platform engineer
- [ ] Service started against a real Vault or AWS Secrets Manager instance in the test environment
- [ ] Startup failure verified when a secret is deliberately missing
- [ ] Log output confirmed to contain no secret values
- [ ] Unit tests pass
- [ ] Story moved to Done in the project tracker

---

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| Epic 1.4 Template Engine | Predecessor | Required |
| Vault or AWS Secrets Manager available in test env | Prerequisite | Infrastructure requirement |

---

## Definition of Done

- All acceptance criteria are met
- Code reviewed and approved
- Committed to `main` via approved PR
