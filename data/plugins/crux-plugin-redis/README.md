# crux-plugin-redis

**Tier:** 1 (Official) | **Version:** 1.0.0 | **Phase:** Pilot

Redis integration for crux-generated Go services. Provides a client with health check registration and graceful shutdown.

## Questions

| ID | Type | Prompt | Default |
|---|---|---|---|
| `redis_distributed_lock` | confirm | Enable Redlock | `false` |
| `redis_ttl_strategy` | select | Default TTL strategy | `fixed` |

## Generated Files

| File | Description |
|---|---|
| `internal/infrastructure/redis/redis.go` | Redis client with structured logging |
| `internal/infrastructure/redis/health.go` | Health check registration |

## Usage in Generated Service

```go
rdb, err := redis.New(ctx, redis.Config{
    Addr:     os.Getenv("REDIS_ADDR"),
    Password: os.Getenv("REDIS_PASSWORD"),
}, log)
rdb.RegisterHealth(registry)
runner.Register(func(ctx context.Context) error { return rdb.Close() })
```

## Environment Variables

| Variable | Description |
|---|---|
| `REDIS_ADDR` | Redis address (`host:6379`) |
| `REDIS_PASSWORD` | Redis password (empty = no auth) |
