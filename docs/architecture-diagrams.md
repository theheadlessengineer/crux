# Crux — Architecture Diagrams

All diagrams reflect the current codebase as of Phase 2 (MVP).

---

## 1. High-Level System Architecture

> What crux is and what it produces.

```mermaid
graph TB
    subgraph ENGINEER["Engineer's Machine"]
        CLI["crux CLI binary\ncmd/crux/main.go"]
    end

    subgraph CRUX["crux Core (internal/)"]
        PRES["presentation/cli\nCobra commands"]
        TUI["presentation/tui\nBubbletea TUI"]
        PROMPT["domain/prompt\nDecision graph + session"]
        PLUGIN_D["domain/plugin\nManifest types + Loader interface"]
        TMPL_D["domain/template\nEngine interface + TemplateData"]
        GEN["infrastructure/generator\nGenerate() — wires everything"]
        TMPL_I["infrastructure/template\nGo text/template engine"]
        PLUGIN_I["infrastructure/plugin\nLoader — embedded + ~/.crux"]
        LOCK["domain/lockfile\n.skeleton.json + crux.lock"]
    end

    subgraph DATA["Embedded Data (data/)"]
        TEMPLATES["data/templates/\ngo-gin / python-fastapi\njava-spring / node-express"]
        PLUGINS["data/plugins/\n24 bundled plugins"]
    end

    subgraph OUTPUT["Generated Service (outputDir/)"]
        APP["Application code\n(language-specific)"]
        TIER1["Tier 1 standards\nresilience.yaml · slo.yaml\nkubernetes/ · monitoring/\ncompliance/ · docs/"]
        PLUGOUT["Plugin code\npostgres.go · redis.go\nkafka.go · jwt.go · …"]
        LOCKOUT[".skeleton.json\ncrux.lock"]
    end

    CLI --> PRES
    PRES --> TUI
    PRES --> PROMPT
    PROMPT --> PLUGIN_I
    PLUGIN_I --> PLUGINS
    PRES --> GEN
    GEN --> TMPL_I
    GEN --> LOCK
    TMPL_I --> TEMPLATES
    TMPL_I --> PLUGINS
    GEN --> OUTPUT
    LOCK --> LOCKOUT
```

---

## 2. Hexagonal Architecture (Ports & Adapters)

> How the crux codebase itself is structured internally.

```mermaid
graph TB
    subgraph DOMAIN["Domain Layer (no external deps)"]
        DH["domain/health\nRegistry"]
        DM["domain/model\nValidateServiceName"]
        DP["domain/plugin\nManifest · Plugin · Loader"]
        DPR["domain/prompt\nQuestion · DecisionGraph · Session"]
        DT["domain/template\nEngine interface · TemplateData"]
        DL["domain/lockfile\nSkeleton · Write"]
    end

    subgraph INFRA["Infrastructure Layer (adapters)"]
        IG["infrastructure/generator\nGenerate()"]
        IT["infrastructure/template\nengine — text/template"]
        IP["infrastructure/plugin\nloader — fs.FS + ~/.crux"]
        IL["infrastructure/logging\nslog JSON handler"]
        ITR["infrastructure/tracing\nOTel SDK + Gin middleware"]
        IS["infrastructure/shutdown\nSIGTERM runner"]
        IR["infrastructure/resilience\nresilience.yaml parser"]
        ISC["infrastructure/secrets\nVault + AWS SM stubs"]
        IC["infrastructure/config\nYAML config loader"]
    end

    subgraph PRES["Presentation Layer"]
        PC["presentation/cli\nCobra — new · version · system · validate"]
        PT["presentation/tui\nBubbletea — 6 themes · review screen"]
        PH["presentation/http\nGin — health · errors · security · tracing"]
    end

    PRES -->|uses interfaces| DOMAIN
    INFRA -->|implements interfaces| DOMAIN
    PRES -->|calls| INFRA
```

---

## 3. `crux new` — Full Sequence Diagram

> Every step from the engineer running the command to files on disk.

```mermaid
sequenceDiagram
    actor Eng as Engineer
    participant CLI as presentation/cli
    participant TUI as presentation/tui
    participant PL as infrastructure/plugin (LoadFromFS)
    participant PR as domain/prompt (DecisionGraph)
    participant GEN as infrastructure/generator
    participant TE as infrastructure/template (engine)
    participant FS as File System

    Eng->>CLI: crux new payment-service
    CLI->>CLI: ValidateServiceName()
    CLI->>PL: LoadFromFS(data/plugins.FS, version)
    PL-->>CLI: []*Plugin (24 plugins)

    CLI->>PR: NewDecisionGraph(coreQuestions + pluginQuestions)
    PR-->>CLI: graph (DAG, cycle-checked)

    CLI->>TUI: Run(session, serviceName)
    Note over TUI: language · team · module · SLO<br/>_plugins multiselect<br/>per-plugin questions (DependsOn gated)

    loop Each question
        TUI->>PR: session.NextQuestion()
        PR-->>TUI: *Question (visible per DAG)
        TUI->>PR: session.Record(q, answer)
    end

    TUI->>TUI: Review screen — confirm / edit / abort
    TUI-->>CLI: answers map + selectedPlugins

    CLI->>GEN: Generate(ctx, Config{Language, Plugins, Answers}, outputDir)

    GEN->>TE: infratemplate.New() — load data/templates FS
    TE-->>GEN: engine

    loop Core templates (language-specific + shared)
        GEN->>TE: Render(tmplName, data, outPath)
        TE->>FS: write file
    end

    loop Selected plugins
        GEN->>TE: AddFromFS(plugin/templates sub-FS)
        loop Plugin templates (language-resolved)
            GEN->>TE: Render(tmplPath, data, outPath)
            TE->>FS: write file
        end
    end

    GEN->>FS: mkGitkeep (empty dir stubs)

    CLI->>FS: lockfile.Write(.skeleton.json + crux.lock)
    CLI-->>Eng: ✔ skeleton generated · Next steps: cd payment-service && make dev
```

---

## 4. Plugin System — Component Diagram

> How plugins are discovered, loaded, and rendered.

```mermaid
graph TB
    subgraph EMBED["Embedded at compile time"]
        EFS["data/plugins embed.FS\n24 plugin directories"]
    end

    subgraph USER["User-installed (runtime)"]
        HOME["~/.crux/plugins/\n(optional)"]
    end

    subgraph LOADER["infrastructure/plugin/loader.go"]
        LFS["LoadFromFS(fs.FS, version)"]
        LO["loader.Load(version)"]
        VAL["ValidateManifest()"]
        COMPAT["checkCompatibility()\nsemver constraint check"]
    end

    subgraph DOMAIN["domain/plugin"]
        MANIFEST["Manifest\napiVersion · kind · metadata · spec"]
        SPEC["Spec\nquestions · templates\ntemplates_by_language\nhooks · dependencies"]
        PLUGIN["Plugin\nManifest + Hook funcs"]
    end

    subgraph GEN["infrastructure/generator"]
        SEL["SelectedPlugin\nName + Templates (lang-resolved)"]
        RENDER["renderPlugins()\nAddFromFS → Render per template"]
    end

    EFS --> LFS
    HOME --> LO
    LFS --> VAL
    LO --> VAL
    VAL --> COMPAT
    COMPAT --> PLUGIN
    PLUGIN --> MANIFEST
    MANIFEST --> SPEC
    SPEC -->|"TemplatesForLang(language)"| SEL
    SEL --> RENDER
```

---

## 5. Plugin Trust Tiers

```mermaid
graph TD
    subgraph T1["Tier 1 — Official (trustTier: 1)"]
        T1A["Bundled in binary via embed.FS"]
        T1B["Platform Engineering authored"]
        T1C["postgresql · redis · kafka · auth-jwt\nkubernetes · terraform-aws · github-actions\nprometheus · claude-code"]
    end

    subgraph T2["Tier 2 — Verified Community (trustTier: 2)"]
        T2A["Phase 2 plugins — bundled but community-reviewed"]
        T2B["resilience · spiffe · multitenant · mysql\nmongodb · rabbitmq · gitlab-ci · datadog\nterraform-gcp · terraform-azure · grpc\nclaude-api · openai · github-copilot · cursor"]
    end

    subgraph T3["Tier 3 — Unvetted (trustTier: 3)"]
        T3A["Git URL install only"]
        T3B["Not in official registry"]
        T3C["Flagged in crux audit output"]
    end

    subgraph INSTALL["Install path"]
        I1["Tier 1/2 → LoadFromFS(embed.FS)"]
        I2["Tier 3 → ~/.crux/plugins/ → loader.Load()"]
    end

    T1 --> I1
    T2 --> I1
    T3 --> I2
```

---

## 6. Decision Graph — Question Dependency Flow

> How the prompt engine resolves which questions to show.

```mermaid
flowchart TD
    L["language\ngo · python · java · node"]
    T["team"]
    M["module (Go only)"]
    SLO_A["slo_availability"]
    SLO_P["slo_p99_latency_ms"]
    PL["_plugins\nmultiselect"]

    PG_V["crux-plugin-postgresql.pg_version\nDependsOn: _plugins contains postgresql"]
    PG_M_GO["pg_migration_tool.go\nDependsOn: postgresql + language=go"]
    PG_M_PY["pg_migration_tool.python\nDependsOn: postgresql + language=python"]
    PG_M_JV["pg_migration_tool.java\nDependsOn: postgresql + language=java"]
    PG_M_ND["pg_migration_tool.node\nDependsOn: postgresql + language=node"]

    RD["crux-plugin-redis.redis_distributed_lock\nDependsOn: _plugins contains redis"]
    KF["crux-plugin-kafka.kafka_direction\nDependsOn: _plugins contains kafka"]

    L --> M
    L --> PL
    T --> PL
    M --> PL
    SLO_A --> PL
    SLO_P --> PL
    PL --> PG_V
    PL --> RD
    PL --> KF
    PG_V --> PG_M_GO
    PG_V --> PG_M_PY
    PG_V --> PG_M_JV
    PG_V --> PG_M_ND

```

---

## 7. Template Engine — Rendering Pipeline

```mermaid
flowchart LR
    subgraph LOAD["Load phase (New())"]
        EFS2["data/templates embed.FS"]
        WALK["fs.WalkDir — parse all .tmpl"]
        TMPL["*template.Template\n(named by path)"]
    end

    subgraph EXTEND["Extend phase (AddFromFS())"]
        PFS["plugin/templates sub-FS"]
        WALK2["fs.WalkDir — parse plugin .tmpl"]
        MERGE["Merge into same *template.Template"]
    end

    subgraph RENDER["Render phase (Render())"]
        LOOKUP["tmpl.Lookup(templateName)"]
        DATA["TemplateData.ToMap()\nservice · company · resilience\nslo · cost · infra · meta · answers"]
        EXEC["t.Execute(file, data)"]
        OUT["Output file\n(0644 or 0755 for .sh)"]
    end

    EFS2 --> WALK --> TMPL
    PFS --> WALK2 --> MERGE --> TMPL
    TMPL --> LOOKUP --> EXEC
    DATA --> EXEC --> OUT
```

---

## 8. Language × Plugin Template Resolution

> How the generator picks the right template for the selected language.

```mermaid
flowchart TD
    LANG{"cfg.Language"}

    subgraph CORE["Core templates"]
        GO["go-gin/ templates\n(Go + Gin)"]
        PY["python-fastapi/ templates\n(Python + FastAPI)"]
        JV["java-spring/ templates\n(Java + Spring Boot)"]
        ND["node-express/ templates\n(Node.js + Express)"]
        SH["Shared templates\nresilience · slo · k8s · monitoring\ncompliance · docs"]
    end

    subgraph PLUGIN_RES["Plugin template resolution\nSpec.TemplatesForLang(language)"]
        TBL["templates_by_language[language]\n→ language-specific list"]
        TFB["fallback: templates[]\n→ generic list (Go)"]
    end

    subgraph EXAMPLES["Example: crux-plugin-postgresql"]
        EGO["go → postgres.go.tmpl\nhealth.go.tmpl"]
        EPY["python → app/db/postgres.py.tmpl"]
        EJV["java → PostgresConfig.java.tmpl"]
        END["node → src/db/postgres.js.tmpl"]
    end

    LANG -->|go| GO
    LANG -->|python| PY
    LANG -->|java| JV
    LANG -->|node| ND
    LANG --> SH

    LANG --> PLUGIN_RES
    TBL -->|found| EGO
    TBL -->|found| EPY
    TBL -->|found| EJV
    TBL -->|found| END
    TBL -->|not found| TFB
```

---

## 9. TUI State Machine

> States the Bubbletea TUI moves through during `crux new`.

```mermaid
stateDiagram-v2
    [*] --> AskingQuestion : session.NextQuestion() != nil

    AskingQuestion --> AskingQuestion : answer recorded\nnext question visible (DAG)
    AskingQuestion --> NavigatingBack : user presses b
    NavigatingBack --> AskingQuestion : session.Back() ok

    AskingQuestion --> QuitConfirm : ctrl+c
    QuitConfirm --> AskingQuestion : any key (cancel)
    QuitConfirm --> Aborted : y

    AskingQuestion --> ReviewScreen : all questions answered

    ReviewScreen --> Generating : "Confirm & create service"
    ReviewScreen --> AskingQuestion : "Change plugin selection"\nnavigates back to _plugins
    ReviewScreen --> EditPick : "Edit an answer"
    EditPick --> AskingQuestion : select question to re-answer
    EditPick --> ReviewScreen : esc / b

    ReviewScreen --> QuitConfirm : ctrl+c

    Generating --> [*] : files written · lockfile written
    Aborted --> [*] : exit 0 · no files written
```

---

## 10. Generated Service — Tier 1 Standards Applied

> What every generated service gets regardless of language or plugin selection.

```mermaid
mindmap
  root((Generated Service))
    Observability
      Health endpoints
        /health /ready /live /metrics /version
      Structured JSON logging
        slog · traceId · correlationId · service · env
      W3C traceparent propagation
        OTel SDK · Gin middleware · HTTP transport
      Grafana dashboard
        Four golden signals
      Alerting rules
        ServiceDown · HighErrorRate · HighP99Latency
    Resilience
      resilience.yaml
        timeouts · circuit breaker · retry · bulkhead
      slo.yaml
        availability target · p99 latency · error budget
    Security
      Security headers
        CSP · X-Frame-Options · HSTS · Referrer-Policy
      Input sanitisation
        path traversal · null bytes · CORS
      Non-root Dockerfile
        multi-stage · USER nonroot · HEALTHCHECK
      K8s network policies
        default-deny ingress + egress
    Compliance
      cost-budget.yaml
      data-classification.yaml
      log-retention.yaml
      catalog-entry.yaml
    Operations
      Graceful shutdown
        SIGTERM/SIGINT · drain timeout
      Secret rotation watcher
        Vault / AWS SM stub
      CI/CD pipeline
        lint → test → build → SBOM → scan
    Documentation
      runbook.md
      capacity-model.md
      TODO.md
      ADR-001
      CHANGELOG.md
```

---

## 11. Delivery Phases — Timeline

```mermaid
timeline
    title Crux Delivery Phases

    section Phase 0 — Foundation
        Complete : Go module · CI/CD · linting
                 : Pre-commit hooks · test infra

    section Phase 1 — Pilot
        Complete : crux new — full prompt + TUI
                 : Go + Gin templates (44 files)
                 : 9 pilot plugins bundled
                 : .skeleton.json + crux.lock
                 : Tier 1 standards (20 of 38)

    section Phase 2 — MVP (In Progress)
        Done : Python · Java · Node templates
             : 15 additional plugins (manifests)
             : Plugin template rendering wired (US-2022)
             : Bubbletea TUI — 6 themes · review screen
        To Do : crux plugin search/install commands
              : crux audit command
              : Shared component library v1
              : Soft mandate rollout

    section Phase 3 — Core
        Planned : crux upgrade command
                : Full compliance profiles SOC2/GDPR/HIPAA
                : SPIFFE + Resilience plugins
                : Shared CI workflows
                : Platform Terraform modules
                : Hard mandate

    section Phase 4 — Platform Maturity
        Future : IDP self-service portal
               : Backstage integration
               : Custom template support
               : Plugin SDK
```

---

## 12. `crux new` — Combination Validation Rules

```mermaid
flowchart TD
    ANSWERS["Collected answers\n+ selected plugins"]

    ANSWERS --> V1{Outbox pattern\nselected?}
    V1 -->|yes, no DB plugin| ERR1["❌ Invalid: outbox requires a database"]
    V1 -->|yes, DB present| OK1["✅ Valid"]

    ANSWERS --> V2{kafka_dlq = true?}
    V2 -->|kafka not selected| ERR2["❌ Invalid: DLQ requires Kafka"]
    V2 -->|kafka selected| OK2["✅ Valid"]

    ANSWERS --> V3{auth = none?}
    V3 -->|yes| WARN1["⚠️ Warning: no authentication configured"]
    V3 -->|no| OK3["✅ Valid"]

    ANSWERS --> V4{Kafka + PostgreSQL\nboth selected?}
    V4 -->|yes| AUTO1["⚡ Auto-suggest: Outbox pattern recommended"]
    V4 -->|no| OK4["no suggestion"]

    ANSWERS --> V5{k8s_service_mesh\nselected?}
    V5 -->|yes| AUTO2["⚡ Auto: mesh_mode=true\nin-app retry/CB disabled in resilience.yaml"]
    V5 -->|no| OK5["mesh_mode=false\nin-app resilience active"]
```

---

## 13. Plugin Manifest Schema

```mermaid
classDiagram
    class Manifest {
        +string APIVersion
        +string Kind
        +Metadata Metadata
        +Spec Spec
    }

    class Metadata {
        +string Name
        +string Version
        +string Description
        +string Author
        +TrustTier TrustTier
        +[]string Tags
        +string CruxVersionConstraint
    }

    class Spec {
        +[]QuestionSpec Questions
        +[]string Templates
        +map~string~[]string~ TemplatesByLang
        +HooksSpec Hooks
        +[]string Dependencies
        +TemplatesForLang(language) []string
    }

    class QuestionSpec {
        +string ID
        +string Type
        +string Prompt
        +string Help
        +[]string Options
        +map~string~[]string~ OptionsByLang
        +string Default
        +map~string~string~ DefaultByLang
    }

    class HooksSpec {
        +[]string PreGenerate
        +[]string PostGenerate
    }

    class Plugin {
        +Manifest Manifest
        +[]Hook PreGenerate
        +[]Hook PostGenerate
    }

    Manifest "1" --> "1" Metadata
    Manifest "1" --> "1" Spec
    Spec "1" --> "0..*" QuestionSpec
    Spec "1" --> "1" HooksSpec
    Plugin "1" --> "1" Manifest
```

---

## 14. Generated Service — Package Dependency Graph (Go)

> Dependency flow inside a generated Go + Gin service.

```mermaid
graph BT
    subgraph CMD["cmd/service-name/"]
        MAIN["main.go"]
    end

    subgraph PRES_SVC["internal/presentation/http/"]
        ROUTER["router.go\nsecurity · CORS · tracing · logging"]
        HEALTH_H["health.go"]
        SERVER["server.go"]
    end

    subgraph DOMAIN_SVC["internal/domain/"]
        HEALTH_R["health/registry.go"]
    end

    subgraph INFRA_SVC["internal/infrastructure/"]
        LOG["logging/logger.go\nmiddleware.go"]
        TRACE["tracing/provider.go\nmiddleware.go · httpclient.go"]
        SHUT["shutdown/shutdown.go"]
        ERR["errors/handler.go"]
        SEC["secrets/"]
    end

    subgraph CFG["internal/config/"]
        CONFIG["config.go"]
    end

    MAIN --> CONFIG
    MAIN --> LOG
    MAIN --> TRACE
    MAIN --> HEALTH_R
    MAIN --> ROUTER
    MAIN --> SERVER
    MAIN --> SHUT

    ROUTER --> LOG
    ROUTER --> TRACE
    ROUTER --> HEALTH_H
    ROUTER --> ERR

    HEALTH_H --> HEALTH_R
```

---

## 15. Lockfile Schema

> What `.skeleton.json` records after generation.

```mermaid
erDiagram
    SKELETON {
        string crux_version
        datetime generated_at
        string generated_by
    }

    SKELETON_SERVICE {
        string name
        string language
        string framework
        string team
        string compliance_profile
    }

    PLUGIN_ENTRY {
        string name
        string version
    }

    TIER1_STANDARDS {
        bool enforced
        string[] disabled_standards
    }

    DEVIATION {
        string standard
        string justification
        string remediation_by
    }

    SKELETON ||--|| SKELETON_SERVICE : "service"
    SKELETON ||--|{ PLUGIN_ENTRY : "plugins_used"
    SKELETON ||--|| TIER1_STANDARDS : "tier1_standards"
    SKELETON ||--o{ DEVIATION : "deviations"
    SKELETON ||--o{ ANSWERS : "answers"

    ANSWERS {
        string language
        string team
        string module
        string slo_availability
        int slo_p99_latency_ms
        string[] _plugins
        string pg_version
        string pg_migration_tool
    }
```
