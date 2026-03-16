# crux-plugin-kubernetes

**Tier:** 1 (Official) | **Version:** 1.0.0 | **Phase:** Pilot

Kubernetes manifests for crux-generated services — Deployment, Service, HPA, PDB, and Network Policy.

## Questions

| ID | Type | Prompt | Default |
|---|---|---|---|
| `k8s_deployment_strategy` | select | Deployment strategy | `rolling` |
| `k8s_service_mesh` | select | Service mesh | `none` |
| `k8s_min_replicas` | input | Minimum replicas | `2` |
| `k8s_max_replicas` | input | Maximum replicas (HPA) | `10` |

## Generated Files

| File | Description |
|---|---|
| `kubernetes/deployment.yaml` | Non-root, read-only FS, liveness + readiness probes |
| `kubernetes/service.yaml` | ClusterIP service |
| `kubernetes/hpa.yaml` | CPU + memory autoscaling |
| `kubernetes/pdb.yaml` | Pod Disruption Budget (minAvailable: 1) |
| `kubernetes/networkpolicy.yaml` | Default-deny ingress + egress with DNS allowance |

## Security Defaults

- `runAsNonRoot: true` — container runs as UID 65534
- `readOnlyRootFilesystem: true`
- `allowPrivilegeEscalation: false`
- All Linux capabilities dropped
- Default-deny network policy — declare egress dependencies explicitly
