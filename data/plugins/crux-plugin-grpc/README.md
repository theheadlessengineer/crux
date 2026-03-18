# crux-plugin-grpc

Generates a gRPC server with health check, reflection, and graceful shutdown for `{{ service.name }}`.

## Generated Files

| File | Purpose |
|---|---|
| `internal/infrastructure/grpc/server.go` | gRPC server with health and reflection |
| `internal/infrastructure/grpc/health.go` | HTTP health registry bridge |

## Configuration

| Question | Default | Description |
|---|---|---|
| `grpc_port` | `9090` | gRPC server port |
| `grpc_tls` | `true` | Enable TLS |
| `grpc_reflection` | `true` | Enable server reflection |

## Usage

```go
grpcServer := grpc.New(logger)
// Register your service implementations:
// pb.RegisterPaymentServiceServer(grpcServer.Server(), &paymentService{})
grpcServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
go grpcServer.ListenAndServe()
```
