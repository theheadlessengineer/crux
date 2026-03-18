# crux-plugin-rabbitmq

Generates RabbitMQ integration for `{{ service.name }}`: AMQP client, producer, consumer, and DLQ support.

## Generated Files

| File | Purpose |
|---|---|
| `internal/infrastructure/rabbitmq/rabbitmq.go` | AMQP client with exchange declaration |
| `internal/infrastructure/rabbitmq/health.go` | Health check for the health registry |

## Configuration

| Question | Default | Description |
|---|---|---|
| `rmq_direction` | `both` | producer / consumer / both |
| `rmq_exchange_type` | `topic` | direct / topic / fanout / headers |
| `rmq_dlq` | `true` | Enable Dead Letter Queue |

## Usage

```go
client, err := rabbitmq.New(ctx, rabbitmq.Config{
    URL:          os.Getenv("RABBITMQ_URL"),
    ExchangeName: "{{ service.name }}.events",
    ExchangeType: "topic",
    DLQEnabled:   true,
}, logger)
err = client.Publish(ctx, "order.created", payload)
```
