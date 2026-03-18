# crux-plugin-mongodb

Generates MongoDB integration for `{{ service.name }}`: connection pool, health check, and collection helpers.

## Generated Files

| File | Purpose |
|---|---|
| `internal/infrastructure/mongodb/mongodb.go` | MongoDB client with connection pooling |
| `internal/infrastructure/mongodb/health.go` | Health check for the health registry |

## Configuration

| Question | Default | Description |
|---|---|---|
| `mongo_database` | `{{ service.name }}` | Default database name |
| `mongo_tls` | `true` | Enforce TLS |
| `mongo_auth_mechanism` | `SCRAM-SHA-256` | Authentication mechanism |

## Usage

```go
client, err := mongodb.New(ctx, mongodb.Config{
    URI:         os.Getenv("MONGO_URI"),
    Database:    os.Getenv("MONGO_DATABASE"),
    MaxPoolSize: 10,
}, logger)
name, check := client.HealthCheck()
healthRegistry.Register(name, check)
```
