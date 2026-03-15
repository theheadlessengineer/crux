# US-1114 — Generate Grafana Dashboard Stub for Four Golden Signals

**Epic:** 1.1 Tier 1 Standards Generation
**Phase:** 1 — Pilot
**Priority:** Must Have
**Status:** To Do

---

## User Story

As a user of Crux,
I want every generated service to include a Grafana dashboard JSON stub pre-configured for the four golden signals,
so that the service has a working observability dashboard from the first deployment without manual dashboard creation.

---

## Pre-Development Checklist

- [ ] The company's Grafana data source name is agreed (used in the dashboard JSON)
- [ ] The four golden signal metrics names from the generated service are confirmed (from US-1101 metrics endpoint)
- [ ] Grafana version in use is agreed (dashboard JSON schema varies by version)
- [ ] Epic 1.4 Template Engine is in progress or complete
- [ ] Story estimated and accepted into the sprint

---

## Scope

Generate a `dashboard.json` file in the `infra/monitoring/` directory that, when imported into Grafana, displays the four golden signals for the generated service.

### Required Panels

| Panel | Metric |
|---|---|
| Request Rate | `http_requests_total` — requests per second |
| Error Rate | `http_requests_total{status=~"5.."}` / total |
| Latency (p50, p95, p99) | `http_request_duration_seconds` histogram quantiles |
| Saturation (CPU, Memory) | Container CPU and memory usage |

### In Scope

- `dashboard.json` as a Grafana dashboard definition
- Service name substituted via template variable from the service name entered during `crux new`
- Dashboard UID generated from the service name
- Time range default: last 1 hour
- Refresh interval: 30 seconds

### Out of Scope

- Infrastructure-level dashboards (cluster, node)
- Custom business metric panels (service teams add these)
- Dashboard provisioning automation (infrastructure concern)

---

## Acceptance Criteria

- [ ] `dashboard.json` is present in every generated service under `infra/monitoring/`
- [ ] Dashboard imports successfully into Grafana without errors
- [ ] All four golden signal panels display data when the service is running
- [ ] Service name is correctly substituted in the dashboard title and panel queries
- [ ] Dashboard UID is unique per service name

---

## Post-Completion Checklist

- [ ] Code reviewed by at least one other platform engineer
- [ ] Dashboard imported into a test Grafana instance and panels verified
- [ ] Story moved to Done in the project tracker

---

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| Epic 1.4 Template Engine | Predecessor | Required |
| US-1101 Health endpoints and metrics | Predecessor | Metric names must be confirmed |

---

## Definition of Done

- All acceptance criteria are met
- Dashboard verified in Grafana
- Code reviewed and approved
- Committed to `main` via approved PR
