# crux-plugin-multitenant

Generates multi-tenancy context propagation for `{{ service.name }}`: tenant extraction middleware, context helpers, and cache key namespacing.

## Generated Files

| File | Purpose |
|---|---|
| `internal/infrastructure/tenancy/context.go` | Tenant context helpers and cache key namespacing |
| `internal/infrastructure/tenancy/middleware.go` | Gin middleware — extracts tenant from header |

## Configuration

| Question | Default | Description |
|---|---|---|
| `mt_isolation_model` | `row` | schema / row / silo |
| `mt_tenant_header` | `X-Tenant-ID` | HTTP header carrying tenant ID |
| `mt_rate_limit_rps` | `100` | Per-tenant rate limit (req/s) |

## Usage

```go
// Register middleware
router.Use(tenancy.Middleware())

// In handlers
tenantID, err := tenancy.FromContext(ctx)

// Cache keys
key := tenancy.CacheKey(tenantID, "user:123")
```
