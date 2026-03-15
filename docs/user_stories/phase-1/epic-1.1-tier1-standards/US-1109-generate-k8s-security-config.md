# US-1109 — Generate Read-Only Root Filesystem and Network Policy Kubernetes Config

**Epic:** 1.1 Tier 1 Standards Generation
**Phase:** 1 — Pilot
**Priority:** Must Have
**Status:** Done

---

## User Story

As a user of Crux,
I want every generated service to include Kubernetes security configuration with a read-only root filesystem and default-deny network policies,
so that the container cannot write to disk unexpectedly and can only communicate over explicitly permitted network paths.

---

## Pre-Development Checklist

- [ ] Epic 1.4 Template Engine is in progress or complete
- [ ] The company's Kubernetes namespace convention is agreed
- [ ] The default network policy posture is agreed: default-deny ingress and egress, with explicit allow rules
- [ ] Story estimated and accepted into the sprint

---

## Scope

Generate Kubernetes manifests that enforce a read-only root filesystem via `securityContext` and apply default-deny NetworkPolicy resources for both ingress and egress.

### In Scope

- `securityContext` in the Deployment manifest with `readOnlyRootFilesystem: true`
- An `emptyDir` volume mounted at `/tmp` for any runtime write requirements
- `runAsNonRoot: true` and `runAsUser: 65534` in the security context
- `allowPrivilegeEscalation: false`
- `capabilities.drop: ["ALL"]`
- A default-deny ingress NetworkPolicy
- A default-deny egress NetworkPolicy
- An explicit allow rule for egress to DNS (port 53)
- Placeholder comments for service-specific allow rules

### Out of Scope

- Service-specific network policy rules (service teams add these using the provided stubs)
- mTLS policy (Epic 3.3)
- Pod Security Admission configuration (cluster-level, not generated per service)

---

## Acceptance Criteria

- [x] Generated Deployment includes `readOnlyRootFilesystem: true` — `TestDeployment_ReadOnlyRootFilesystem` passes
- [x] Generated Deployment mounts an `emptyDir` at `/tmp` — `TestDeployment_TmpEmptyDirVolume` passes
- [x] Generated Deployment runs as UID 65534 with `runAsNonRoot: true` — `TestDeployment_RunsAsNonRoot` passes; matches UID 65534 from US-1108 Dockerfile
- [x] `allowPrivilegeEscalation: false` is set — `TestDeployment_AllowPrivilegeEscalationFalse` passes
- [x] All Linux capabilities are dropped — `capabilities.drop: ["ALL"]`; `TestDeployment_DropsAllCapabilities` passes
- [x] Default-deny ingress NetworkPolicy is generated — `ingress: []` with `policyTypes: [Ingress]`; `TestNetworkPolicyIngress_DefaultDeny` passes
- [x] Default-deny egress NetworkPolicy is generated — `policyTypes: [Egress]` with only DNS allow rule; `TestNetworkPolicyEgress_DefaultDeny` passes
- [x] Egress to DNS (port 53) is explicitly allowed — UDP and TCP port 53; `TestNetworkPolicyEgress_DNSAllowed` passes
- [ ] `kubectl apply --dry-run=client` succeeds on all generated manifests — requires a running cluster or kubectl; manual verification step
- [ ] Generated service starts successfully with `readOnlyRootFilesystem: true` active — requires a running cluster; manual verification step

---

## Post-Completion Checklist

- [ ] Code reviewed by at least one other platform engineer
- [ ] Manifests applied to a test cluster and service started successfully
- [ ] Attempt to write to a non-`/tmp` path confirmed to fail
- [ ] Network policy verified: traffic not matching allow rules is dropped
- [x] Unit tests pass — `go test ./data/templates/...` all green (18 tests); `golangci-lint` 0 issues
- [ ] Story moved to Done in the project tracker

---

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| Epic 1.4 Template Engine | Predecessor | Required |
| US-1108 Non-root Dockerfile | Predecessor | UID must match between Dockerfile and K8s config |

---

## Definition of Done

- All acceptance criteria are met
- Code reviewed and approved
- Committed to `main` via approved PR
