# crux-plugin-resilience

Generates resilience patterns for `{{ service.name }}`: circuit breaker, bulkhead, timeout hierarchy, and retry budget.

## Generated Files

| File | Purpose |
|---|---|
| `internal/infrastructure/resilience/circuitbreaker.go` | Circuit breaker with configurable failure threshold |

## Configuration

Answers to plugin questions are written into `resilience.yaml` at the service root.

| Question | Default | Description |
|---|---|---|
| `resilience_profile` | `standard` | Preset: standard / strict / minimal |
| `resilience_cb_threshold` | `50` | % failures before circuit opens |
| `resilience_timeout_http_ms` | `5000` | Outbound HTTP timeout in ms |

## Usage

```go
cb := resilience.NewCircuitBreaker(resilience.DefaultConfig(), logger)
err := cb.Execute(ctx, func(ctx context.Context) error {
    return callDownstream(ctx)
})
```
