# syntax=docker/dockerfile:1

# ── Stage 1: Build ────────────────────────────────────────────────────────────
FROM golang:1.26 AS builder

WORKDIR /app

# Download dependencies first — cached unless go.mod/go.sum change.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION=dev
# CGO_ENABLED=0 required for distroless (no C runtime in final image).
RUN CGO_ENABLED=0 GOOS=linux go build \
      -ldflags="-X main.version=${VERSION}" \
      -o /service \
      ./cmd/my-service

# ── Stage 2: Final ────────────────────────────────────────────────────────────
# distroless/static:nonroot — no shell, no package manager, runs as UID 65534.
FROM gcr.io/distroless/static:nonroot

COPY --from=builder /service /service

USER nonroot:nonroot

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD ["/service", "healthcheck"]

ENTRYPOINT ["/service"]
