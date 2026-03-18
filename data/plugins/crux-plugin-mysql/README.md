# crux-plugin-mysql

Generates MySQL integration for `{{ service.name }}`: connection pool, health check, and migration tooling stub.

## Generated Files

| File | Purpose |
|---|---|
| `internal/infrastructure/mysql/mysql.go` | Connection pool with structured logging |
| `internal/infrastructure/mysql/health.go` | Health check for the health registry |

## Configuration

| Question | Default | Description |
|---|---|---|
| `mysql_version` | `8.0` | MySQL version to target |
| `mysql_migration_tool` | `goose` | Migration tool |
| `mysql_tls` | `true` | Enforce TLS on connections |

## Usage

```go
db, err := mysql.New(ctx, mysql.Config{
    DSN:          os.Getenv("MYSQL_DSN"),
    MaxOpenConns: 10,
    MaxIdleConns: 5,
}, logger)
name, check := db.HealthCheck()
healthRegistry.Register(name, check)
```
