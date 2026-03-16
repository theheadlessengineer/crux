# ADR-003 — Plugin Isolation Strategy

**Date:** 2026-03-16
**Status:** Accepted
**Deciders:** Platform Engineering Team

---

## Context

The plugin system (Epic 1.6) requires a decision on how plugins are isolated from the crux core process and from each other. The options are:

1. **In-process** — plugins are compiled Go packages loaded at startup; no process boundary.
2. **Sub-process** — each plugin runs as a child process; communication over stdin/stdout or gRPC.
3. **WASM sandbox** — plugins compiled to WASM, executed in a sandboxed runtime (e.g. Wazero).

The choice affects security, performance, developer experience, and implementation complexity.

---

## Decision

**In-process with interface constraints** for Phase 1 (Pilot) and Phase 2 (MVP).

Plugins are Go packages that implement the `domain.Plugin` interface. They are compiled into the crux binary or loaded from `~/.crux/plugins/` as pre-compiled binaries at startup. No process boundary exists between the plugin and the core.

---

## Rationale

| Criterion | In-process | Sub-process | WASM |
|---|---|---|---|
| Implementation complexity | Low | High | High |
| Performance | Excellent | Moderate (IPC overhead) | Moderate |
| Developer experience | Excellent (standard Go) | Moderate | Poor (WASM toolchain) |
| Security isolation | None | Strong | Strong |
| Phase 1 feasibility | ✅ | ❌ | ❌ |

At Phase 1, all plugins are Tier 1 (platform-team owned) or Tier 2 (inner-source, reviewed). The risk of malicious plugin code is low and mitigated by the trust tier model and code review process.

---

## Known Limitations

- A buggy plugin can crash the crux process.
- A malicious Tier 3 plugin has full access to the host filesystem and network.
- `postGenerate` hooks run with the same OS permissions as the crux process.

---

## Mitigations

- Tier 3 plugins display an explicit security warning before installation and require user confirmation.
- `crux audit` flags services using Tier 3 plugins.
- Plugin manifests are validated before any hook is executed.
- Version compatibility is checked at startup — incompatible plugins are rejected before hooks run.

---

## Future Path

The `domain.Loader` interface is designed to be replaced. When WASM sandboxing becomes feasible (Phase 3+), a `WASMLoader` can be substituted without changing the core or any plugin consumer. The interface contract is:

```go
type Loader interface {
    Load(cruxVersion string) ([]*Plugin, error)
}
```

Open Decision #3 in the scope document (WASM sandbox) is deferred to Phase 3.

---

## Consequences

- Phase 1 plugin development uses standard Go — no special toolchain required.
- The `plugin.yaml` manifest schema is the stable contract; the execution model may change.
- All Tier 1 and Tier 2 plugins must pass code review before being accepted into the registry.
