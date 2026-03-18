# crux-plugin-spiffe

Generates SPIFFE/SPIRE workload identity wiring for `{{ service.name }}`: SVID management client and OPA policy stub.

## Generated Files

| File | Purpose |
|---|---|
| `internal/infrastructure/identity/workload.go` | SPIRE workload API client — SVID fetch |
| `identity/policies/service.rego` | OPA policy stub — service-to-service authorisation |

## Configuration

| Question | Default | Description |
|---|---|---|
| `spiffe_trust_domain` | `company.com` | SPIFFE trust domain |
| `spiffe_socket_path` | `/run/spire/sockets/agent.sock` | SPIRE agent socket |

## Usage

```go
client, err := identity.New(ctx, identity.Config{
    SocketPath:  os.Getenv("SPIRE_SOCKET"),
    TrustDomain: "company.com",
    ServiceName: "{{ service.name }}",
}, logger)
```
