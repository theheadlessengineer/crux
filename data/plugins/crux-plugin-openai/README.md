# crux-plugin-openai

Generates an OpenAI API client with safety guardrails for `{{ service.name }}`.

## Generated Files

| File | Purpose |
|---|---|
| `internal/infrastructure/ai/openai/client.go` | OpenAI chat completions client with cost circuit breaker |
| `internal/infrastructure/ai/openai/pii_scrubber.go` | PII scrubbing and injection detection |

## Safety Guardrails

| Guardrail | Mechanism |
|---|---|
| PII scrubbing | Regex patterns from `data-classification.yaml` applied before every call |
| Cost circuit breaker | Opens at ${{ openai_daily_budget_usd }}/day — resets at midnight UTC |
| Audit logging | Every call logged with token counts (no payload) |
| Injection detection | `DetectInjection()` scans user input before forwarding |

## Configuration

| Question | Default | Description |
|---|---|---|
| `openai_model` | `gpt-4o` | Default model |
| `openai_daily_budget_usd` | `50` | Daily spend circuit breaker (USD) |
| `openai_max_tokens` | `1024` | Max tokens per response |

## Usage

```go
scrubber := openai.NewPIIScrubber(nil)
client := openai.New(openai.Config{APIKey: os.Getenv("OPENAI_API_KEY")}, scrubber, logger)
response, err := client.Chat(ctx, userPrompt)
```
