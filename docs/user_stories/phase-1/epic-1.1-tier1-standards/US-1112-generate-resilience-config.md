# US-1112 — Generate Resilience Configuration (resilience.yaml)

**Epic:** 1.1 Tier 1 Standards Generation
**Phase:** 1 — Pilot
**Priority:** Must Have
**Status:** To Do

---

## User Story

As a user of Crux,
I want every generated service to include a `resilience.yaml` file with default timeout, retry, circuit breaker, and bulkhead configuration,
so that the service has documented, reviewable resilience settings from day one rather than relying on library defaults.

---

## Pre-Development Checklist

- [ ] The agreed default resilience values are documented (timeout: 5s, retries: 3, circuit breaker threshold: 50%)
- [ ] Epic 1.4 Template Engine is in progress or complete
- [ ] Story estimated and accepted into the sprint

---

## Scope

Generate a `resilience.yaml` configuration file and a configuration loader that applies the declared values to all outbound HTTP clients.

### resilience.yaml Required Sections

- `timeout`: Default call timeout in milliseconds
- `retry`: Max attempts, backoff type (exponential), initial interval, max interval
- `circuitBreaker`: Failure rate threshold, slow call threshold, minimum call count, open state wait duration
- `bulkhead`: Maximum concurrent calls per downstream

### In Scope

- `resilience.yaml` template with sensible defaults
- Configuration loader parsing the YAML into typed structs
- Outbound HTTP client configured with the resilience parameters
- Unit tests verifying configuration loading and default values

### Out of Scope

- Full resilience plugin implementation (Epic 3.2)
- Service mesh integration (Epic 2.3 and Epic 3.2)

---

## Acceptance Criteria

- [ ] `resilience.yaml` is present in every generated service
- [ ] Timeout, retry, circuit breaker, and bulkhead sections are present with documented defaults
- [ ] Configuration loader reads the file and populates typed structs
- [ ] Outbound HTTP client uses the loaded values
- [ ] Unit tests verify default values and override via configuration
- [ ] Invalid configuration causes startup failure with a clear error message

---

## Post-Completion Checklist

- [ ] Code reviewed by at least one other platform engineer
- [ ] Configuration loaded and verified at runtime
- [ ] Unit tests pass
- [ ] Story moved to Done in the project tracker

---

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| Epic 1.4 Template Engine | Predecessor | Required |

---

## Definition of Done

- All acceptance criteria are met
- Code reviewed and approved
- Committed to `main` via approved PR
