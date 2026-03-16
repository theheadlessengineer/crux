# crux-plugin-kafka

**Tier:** 1 (Official) | **Version:** 1.0.0 | **Phase:** Pilot

Kafka integration for crux-generated Go services. Provides a producer/consumer client with DLQ support, health check registration, and graceful shutdown.

## Questions

| ID | Type | Prompt | Default |
|---|---|---|---|
| `kafka_direction` | select | Producer / consumer / both | `both` |
| `kafka_dlq` | confirm | Enable DLQ | `true` |
| `kafka_outbox` | confirm | Enable outbox pattern | `false` |
| `kafka_schema_format` | select | Message schema format | `json` |

## Generated Files

| File | Description |
|---|---|
| `internal/infrastructure/kafka/kafka.go` | Kafka client (franz-go) |
| `internal/infrastructure/kafka/health.go` | Health check registration |

## Usage in Generated Service

```go
kc, err := kafka.New(ctx, kafka.Config{
    Brokers:       strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
    ConsumerGroup: os.Getenv("KAFKA_CONSUMER_GROUP"),
    Topics:        []string{os.Getenv("KAFKA_TOPIC")},
    DLQTopic:      os.Getenv("KAFKA_DLQ_TOPIC"),
}, log)
kc.RegisterHealth(registry)
runner.Register(func(ctx context.Context) error { kc.Close(); return nil })
```

## Environment Variables

| Variable | Description |
|---|---|
| `KAFKA_BROKERS` | Comma-separated broker list |
| `KAFKA_CONSUMER_GROUP` | Consumer group ID |
| `KAFKA_TOPIC` | Primary topic |
| `KAFKA_DLQ_TOPIC` | Dead letter queue topic |
