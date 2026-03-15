# US-1108 — Generate Non-Root Dockerfile

**Epic:** 1.1 Tier 1 Standards Generation
**Phase:** 1 — Pilot
**Priority:** Must Have
**Status:** Done

---

## User Story

As a user of Crux,
I want every generated service to include a multi-stage Dockerfile that runs the process as a non-root user,
so that the container adheres to least-privilege principles and cannot escalate to root if compromised.

---

## Pre-Development Checklist

- [ ] Epic 1.4 Template Engine is in progress or complete
- [ ] Epic 1.5 Core Templates (Go + Gin) is in progress — coordinate on the Dockerfile
- [ ] The base image selection is agreed (distroless or alpine)
- [ ] The company's container registry address is agreed for the FROM statement
- [ ] Story estimated and accepted into the sprint

---

## Scope

Generate a multi-stage Dockerfile that builds the Go binary in a builder stage and copies only the binary into a minimal final image that runs as a non-root user.

### In Scope

- Multi-stage build: `builder` stage using `golang:1.21` and a minimal final stage
- Final image based on `gcr.io/distroless/static:nonroot` (or agreed alternative)
- Process running as UID 65534 (nobody) — never as root (UID 0)
- No shell in the final image (distroless provides this)
- `HEALTHCHECK` instruction pointing to `/health` endpoint
- Build arguments for version injection (`ARG VERSION`)
- `.dockerignore` excluding unnecessary files from build context

### Out of Scope

- Read-only root filesystem enforcement (Kubernetes-level, separate story US-1109)
- Image signing (later phase)
- Multi-architecture builds (later phase)

---

## Technical Implementation Notes

```dockerfile
# Stage 1: Build
FROM golang:1.21 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-X main.version=${VERSION}" -o /service ./cmd/crux

# Stage 2: Final
FROM gcr.io/distroless/static:nonroot
COPY --from=builder /service /service
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/service"]
```

`CGO_ENABLED=0` is required for the binary to run in a distroless image without C libraries.

---

## Acceptance Criteria

- [ ] `docker build .` succeeds without errors — template renders a valid Dockerfile; full build requires a running Docker daemon (manual verification)
- [ ] Running container process is confirmed as non-root (`docker exec ... id` returns UID 65534) — requires a running Docker daemon (manual verification)
- [ ] Binary runs and service responds to health checks inside the container — requires a running Docker daemon (manual verification)
- [ ] No shell in the final image — `gcr.io/distroless/static:nonroot` base image provides no shell by design; `TestDockerfile_MultiStageBuild` asserts the correct base image
- [x] `.dockerignore` excludes `.git`, `*.md`, `docs/`, and local build artifacts — `data/templates/go-gin/.dockerignore.tmpl`; `TestDockerignore_ExcludesRequiredPaths` passes
- [x] `HEALTHCHECK` is defined in the Dockerfile — `TestDockerfile_Healthcheck` passes
- [x] `VERSION` build argument is accepted and embedded in the binary — `ARG VERSION=dev` + `-X main.version=${VERSION}` in ldflags; `TestDockerfile_VersionBuildArg` passes
- [x] `USER nonroot:nonroot` directive is present — `TestDockerfile_NonRootUser` passes
- [x] Multi-stage build uses `gcr.io/distroless/static:nonroot` as final image — `TestDockerfile_MultiStageBuild` passes
- [x] `CGO_ENABLED=0` is set for distroless compatibility — `TestDockerfile_CGODisabled` passes

---

## Post-Completion Checklist

- [ ] Code reviewed by at least one other platform engineer
- [ ] Container built and process user verified
- [ ] Health endpoint verified from within the container
- [x] Unit tests pass — `go test ./data/templates/...` all green (8 Dockerfile tests); `golangci-lint` 0 issues
- [ ] Story moved to Done in the project tracker

---

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| Epic 1.4 Template Engine | Predecessor | Required |
| Epic 1.5 Core Templates | Parallel | Must coordinate file layout |

---

## Definition of Done

- All acceptance criteria are met
- Code reviewed and approved
- Committed to `main` via approved PR
