# crux-plugin-claude-api

Generates an Anthropic Claude API client with safety guardrails for `{{ service.name }}`.

## Generated Files

| File | Purpose |
|---|---|
| `internal/infrastructure/ai/claude/client.go` | Claude API client with cost circuit breaker |
| `internal/infrastructure/ai/claude/pii_scrubber.go` | PII scrubbing and injection detection |

## Safety Guardrails

| Guardrail | Mechanism |
|---|---|
| PII scrubbing | Regex patterns from `data-classification.yaml` applied before every call |
| Cost circuit breaker | Opens at ${{ claude_daily_budget_usd }}/day — resets at midnight UTC |
| Audit logging | Every call logged with token counts (no payload) |
| Injection detection | `DetectInjection()` scans user input before forwarding |

## Configuration

| Question | Default | Description |
|---|---|---|
| `claude_model` | `claude-3-5-sonnet-20241022` | Default model |
| `claude_daily_budget_usd` | `50` | Daily spend circuit breaker (USD) |
| `claude_max_tokens` | `1024` | Max tokens per response |

## Usage

```go
scrubber := claude.NewPIIScrubber(nil) // load patterns from data-classification.yaml
client := claude.New(claude.Config{APIKey: os.Getenv("ANTHROPIC_API_KEY")}, scrubber, logger)
response, err := client.Message(ctx, userPrompt)
```
