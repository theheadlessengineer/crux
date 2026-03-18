# crux-plugin-datadog

Generates Datadog APM, metrics, and log forwarding configuration for `{{ service.name }}`.

## Generated Files

| File | Purpose |
|---|---|
| `internal/infrastructure/datadog/datadog.go` | APM tracer init and shutdown |
| `kubernetes/datadog-agent.yaml` | Kubernetes ConfigMap for DD environment variables |

## Configuration

| Question | Default | Description |
|---|---|---|
| `dd_service_env` | `production` | Datadog env tag |
| `dd_apm_enabled` | `true` | Enable APM tracing |
| `dd_log_injection` | `true` | Inject trace IDs into logs |

## Usage

```go
// In main.go
datadog.Start(datadog.DefaultConfig("{{ service.name }}", version), logger)
defer datadog.Stop()
```
