# crux-plugin-postgresql

**Tier:** 1 (Official) | **Version:** 1.0.0 | **Phase:** Pilot

PostgreSQL integration for crux-generated Go services. Provides a connection pool, health check registration, and graceful shutdown.

## Questions

| ID | Type | Prompt | Default |
|---|---|---|---|
| `pg_version` | select | PostgreSQL version | `16` |
| `pg_read_replica` | confirm | Enable read replica | `false` |
| `pg_audit_log` | select | Audit logging strategy | `application` |
| `pg_migration_tool` | select | Migration tool | `goose` |

## Generated Files

| File | Description |
|---|---|
| `internal/infrastructure/postgres/postgres.go` | Connection pool with structured logging |
| `internal/infrastructure/postgres/health.go` | Health check registration |

## Usage in Generated Service

```go
db, err := postgres.New(ctx, postgres.Config{
    DSN:             os.Getenv("DATABASE_URL"),
    MaxOpenConns:    10,
    MaxIdleConns:    5,
    ConnMaxLifetime: 5 * time.Minute,
}, log)
db.RegisterHealth(registry)
runner.Register(func(ctx context.Context) error { return db.Close() })
```

## Environment Variables

| Variable | Description |
|---|---|
| `DATABASE_URL` | PostgreSQL DSN (`postgres://user:pass@host:5432/db?sslmode=require`) |
