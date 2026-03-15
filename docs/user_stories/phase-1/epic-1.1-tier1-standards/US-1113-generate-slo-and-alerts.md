# US-1113 — Generate SLO Definition and Baseline Alert Rules

**Epic:** 1.1 Tier 1 Standards Generation
**Phase:** 1 — Pilot
**Priority:** Must Have
**Status:** To Do

---

## User Story

As a user of Crux,
I want every generated service to include an `slo.yaml` file and `alerts.yaml` with baseline alerting rules for the four golden signals,
so that the service has measurable availability and latency targets and will alert on violations from the first deployment.

---

## Pre-Development Checklist

- [ ] The company's default SLO targets are agreed (e.g., 99.9% availability, p99 latency < 500ms)
- [ ] The alerting backend is agreed (Prometheus + AlertManager)
- [ ] The four golden signals are agreed: error rate, latency, traffic, saturation
- [ ] Epic 1.4 Template Engine is in progress or complete
- [ ] Story estimated and accepted into the sprint

---

## Scope

Generate `slo.yaml` and `alerts.yaml` files that define baseline service level objectives and the four golden signal alert rules.

### slo.yaml Required Fields

- Service name
- Availability target (e.g., 99.9%)
- Latency target (p99, p95 in milliseconds)
- Error budget window (e.g., 30 days)

### alerts.yaml — Four Baseline Alerts

| Alert | Condition |
|---|---|
| High Error Rate | HTTP 5xx rate > 1% over 5 minutes |
| High Latency | p99 latency > configured target for 5 minutes |
| High Traffic | Request rate > capacity model threshold |
| High Saturation | CPU or memory > 80% for 10 minutes |

### In Scope

- `slo.yaml` template with company default targets as stubs (teams fill in service-specific values)
- `alerts.yaml` with four Prometheus alerting rules
- Placeholder thresholds documented with comments explaining how to calibrate them

### Out of Scope

- SLO burn rate alerts (advanced, Epic 3.4)
- Custom business metric alerts (service teams add these)
- Grafana dashboard (US-1114)

---

## Acceptance Criteria

- [ ] `slo.yaml` is present in every generated service with availability, latency, and error budget sections
- [ ] `alerts.yaml` contains all four golden signal alert rules
- [ ] Alert rule syntax is valid Prometheus YAML (validated with `promtool check rules`)
- [ ] Default thresholds are documented with comments
- [ ] Thresholds reference the values in `slo.yaml` for consistency

---

## Post-Completion Checklist

- [ ] Code reviewed by at least one other platform engineer
- [ ] `promtool check rules alerts.yaml` passes with no errors
- [ ] SLO targets reviewed by the team for reasonableness
- [ ] Story moved to Done in the project tracker

---

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| Epic 1.4 Template Engine | Predecessor | Required |

---

## Definition of Done

- All acceptance criteria are met
- Alert rule syntax validated
- Code reviewed and approved
- Committed to `main` via approved PR
