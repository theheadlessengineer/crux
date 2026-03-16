# crux-plugin-prometheus

**Tier:** 1 (Official) | **Version:** 1.0.0 | **Phase:** Pilot

Prometheus alerting rules and Grafana dashboard (four golden signals) for crux-generated services.

## Questions

| ID | Type | Prompt | Default |
|---|---|---|---|
| `prom_alerting_backend` | select | Alert routing backend | `alertmanager` |
| `prom_p99_latency_ms` | input | p99 latency alert threshold (ms) | `500` |
| `prom_error_rate_threshold` | input | Error rate alert threshold (%) | `5` |

## Generated Files

| File | Description |
|---|---|
| `monitoring/alerts.yaml` | Four baseline alerts (down, error rate, p99 latency, memory) |
| `monitoring/dashboard.json` | Grafana dashboard — request rate, error rate, latency, CPU |

## Alerts

| Alert | Condition | Severity |
|---|---|---|
| `ServiceDown` | No healthy instances for 1m | critical |
| `HighErrorRate` | Error rate > threshold for 5m | warning |
| `HighP99Latency` | p99 > threshold for 5m | warning |
| `HighMemoryUsage` | Memory > 85% of limit for 10m | warning |
