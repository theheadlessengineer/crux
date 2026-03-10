# Plugin Development Guide

## Overview

Plugins extend crux functionality without modifying core code. Each plugin is a self-contained module that contributes questions, templates, validators, and hooks.

## Plugin Architecture

### Plugin Types

**Tier 1 - Official Plugins:**
- Maintained by platform team
- Bundled with crux core
- Full capability access
- Examples: PostgreSQL, Redis, Kafka, Kubernetes

**Tier 2 - Community Plugins:**
- Community maintained
- Reviewed by platform team
- Listed in official registry
- Examples: RabbitMQ, MongoDB, Datadog

**Tier 3 - Unvetted Plugins:**
- Self-published
- Not in official registry
- Sandboxed execution
- Security warnings displayed

### Plugin Structure

```
crux-plugin-postgresql/
├── plugin.yaml              # Plugin manifest
├── questions.yaml           # Interactive prompts
├── templates/               # Template files
│   ├── go/
│   │   ├── database.go.tmpl
│   │   └── repository.go.tmpl
│   ├── python/
│   │   └── database.py.tmpl
│   └── java/
│       └── Database.java.tmpl
├── validators.yaml          # Validation rules
├── hooks/                   # Lifecycle hooks
│   ├── pre-generate.sh
│   └── post-generate.sh
├── tests/                   # Plugin tests
│   └── plugin_test.go
└── README.md
```

## Plugin Manifest

### plugin.yaml

```yaml
name: crux-plugin-postgresql
version: 2.1.0
tier: 1
description: PostgreSQL integration with connection pooling and migrations
author: platform-team@company.com
crux_core_compatibility: ">=1.0.0 <3.0.0"

languages_supported:
  - go
  - python
  - java
  - node

capabilities:
  - questions
  - templates
  - hooks
  - validators

depends_on:
  - crux-plugin-base-db@^1.0.0

questions_schema: ./questions.yaml
templates_dir: ./templates/
hooks_dir: ./hooks/
validators_file: ./validators.yaml
tests_dir: ./tests/
```

### Required Fields

- `name`: Unique plugin identifier (must start with `crux-plugin-`)
- `version`: Semantic version
- `tier`: 1, 2, or 3
- `description`: Brief description
- `author`: Contact email
- `crux_core_compatibility`: Semver range

### Optional Fields

- `languages_supported`: List of supported languages
- `capabilities`: List of plugin capabilities
- `depends_on`: Plugin dependencies
- `homepage`: Plugin documentation URL
- `repository`: Source code URL
- `license`: SPDX license identifier

## Questions Schema

### questions.yaml

```yaml
questions:
  - id: db_read_replica
    prompt: "Enable read replica support?"
    type: confirm
    default: false
    group: "Data Layer"
    order: 20
    depends_on:
      question: db_type
      value: postgresql
    help: "Configures read-only replica for query offloading"

  - id: db_pool_size
    prompt: "Connection pool size?"
    type: number
    default: 10
    min: 1
    max: 100
    group: "Data Layer"
    order: 21

  - id: db_migration_tool
    prompt: "Migration tool?"
    type: select
    options:
      - label: "Alembic (Python)"
        value: alembic
      - label: "Flyway (Java)"
        value: flyway
      - label: "golang-migrate"
        value: golang-migrate
    default: alembic
    group: "Data Layer"
    order: 22
```

### Question Types

**confirm:**
```yaml
type: confirm
default: false
```

**text:**
```yaml
type: text
default: "localhost"
validation: "^[a-z0-9.-]+$"
```

**number:**
```yaml
type: number
default: 10
min: 1
max: 100
```

**select:**
```yaml
type: select
options:
  - label: "Option 1"
    value: opt1
  - label: "Option 2"
    value: opt2
default: opt1
```

**multiselect:**
```yaml
type: multiselect
options:
  - label: "Feature A"
    value: feature_a
  - label: "Feature B"
    value: feature_b
default: [feature_a]
```

### Conditional Questions

```yaml
depends_on:
  question: db_type
  value: postgresql

# OR multiple conditions
depends_on:
  all:
    - question: db_type
      value: postgresql
    - question: environment
      value: production
```

## Templates

### Template Syntax

Templates use Go's `text/template` syntax:

```go
// database.go.tmpl
package database

import (
    "database/sql"
    _ "github.com/lib/pq"
)

type Config struct {
    Host     string
    Port     int
    Database string
    User     string
    Password string
    {{- if .db_read_replica }}
    ReplicaHost string
    {{- end }}
    PoolSize int
}

func NewConnection(cfg *Config) (*sql.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s port=%d dbname=%s user=%s password=%s pool_max_conns=%d",
        cfg.Host, cfg.Port, cfg.Database, cfg.User, cfg.Password, cfg.PoolSize,
    )
    
    return sql.Open("postgres", dsn)
}
```

### Available Variables

**Service Variables:**
```
{{ .service.name }}           # payment-service
{{ .service.language }}       # go
{{ .service.framework }}      # gin
```

**Plugin Variables:**
```
{{ .db_read_replica }}        # true/false
{{ .db_pool_size }}           # 10
{{ .db_migration_tool }}      # alembic
```

**Company Variables:**
```
{{ .company.name }}           # Acme Corp
{{ .company.registry }}       # registry.company.com
```

### Template Functions

**String Functions:**
```
{{ .service.name | upper }}              # PAYMENT-SERVICE
{{ .service.name | lower }}              # payment-service
{{ .service.name | title }}              # Payment-Service
{{ .service.name | replace "-" "_" }}    # payment_service
```

**Conditional Functions:**
```
{{ if .db_read_replica }}
// Read replica code
{{ end }}

{{ if not .db_read_replica }}
// Single instance code
{{ end }}
```

**Loop Functions:**
```
{{ range .environments }}
- {{ . }}
{{ end }}
```

### Multi-Language Templates

Organize by language:

```
templates/
├── go/
│   ├── database.go.tmpl
│   └── repository.go.tmpl
├── python/
│   ├── database.py.tmpl
│   └── repository.py.tmpl
└── java/
    └── Database.java.tmpl
```

Template selection based on `{{ .service.language }}`.

## Validators

### validators.yaml

```yaml
validators:
  - name: read_replica_requires_primary
    condition: "db_read_replica == true && db_type == ''"
    error: "Read replica requires primary database configuration"
    severity: error

  - name: pool_size_warning
    condition: "db_pool_size > 50"
    error: "Pool size > 50 may cause resource exhaustion"
    severity: warning

  - name: migration_tool_compatibility
    condition: "db_migration_tool == 'alembic' && service.language != 'python'"
    error: "Alembic requires Python"
    severity: error
```

### Validator Fields

- `name`: Unique validator identifier
- `condition`: Boolean expression
- `error`: Error message
- `severity`: `error` (blocks generation) or `warning` (shows warning)

### Expression Syntax

```yaml
# Equality
condition: "db_type == 'postgresql'"

# Inequality
condition: "db_pool_size > 50"

# Logical operators
condition: "db_read_replica == true && db_type == 'postgresql'"

# Negation
condition: "!(db_type == 'mysql' || db_type == 'postgresql')"
```

## Hooks

### Lifecycle Hooks

**pre-generate.sh:**
Runs before template rendering.

```bash
#!/bin/bash
set -e

# Validate PostgreSQL is available
if ! command -v psql &> /dev/null; then
    echo "Warning: psql not found. Install PostgreSQL client."
fi

# Check network connectivity
if ! nc -z localhost 5432; then
    echo "Warning: PostgreSQL not running on localhost:5432"
fi
```

**post-generate.sh:**
Runs after template rendering.

```bash
#!/bin/bash
set -e

SERVICE_DIR=$1

# Install dependencies
cd "$SERVICE_DIR"

case "$LANGUAGE" in
    go)
        go get github.com/lib/pq
        go mod tidy
        ;;
    python)
        pip install psycopg2-binary alembic
        ;;
    java)
        # Maven dependencies already in pom.xml
        ;;
esac

# Initialize migration directory
if [ "$DB_MIGRATION_TOOL" = "alembic" ]; then
    alembic init migrations
fi

echo "PostgreSQL plugin setup complete"
```

### Hook Environment Variables

Available in all hooks:

```bash
CRUX_VERSION=1.0.0
SERVICE_NAME=payment-service
SERVICE_LANGUAGE=go
SERVICE_FRAMEWORK=gin
OUTPUT_DIR=/path/to/output
PLUGIN_NAME=crux-plugin-postgresql
PLUGIN_VERSION=2.1.0

# All question answers as environment variables
DB_READ_REPLICA=true
DB_POOL_SIZE=10
DB_MIGRATION_TOOL=alembic
```

### Hook Exit Codes

- `0`: Success
- `1`: Error (blocks generation)
- `2`: Warning (shows warning, continues)

## Testing Plugins

### Test Structure

```go
// tests/plugin_test.go
package tests

import (
    "testing"
    "os"
    "path/filepath"
)

func TestPluginManifest(t *testing.T) {
    manifest, err := LoadManifest("../plugin.yaml")
    if err != nil {
        t.Fatalf("failed to load manifest: %v", err)
    }
    
    if manifest.Name != "crux-plugin-postgresql" {
        t.Errorf("unexpected name: %s", manifest.Name)
    }
}

func TestTemplateRendering(t *testing.T) {
    tests := []struct {
        name     string
        language string
        answers  map[string]interface{}
    }{
        {
            name:     "go with read replica",
            language: "go",
            answers: map[string]interface{}{
                "db_read_replica": true,
                "db_pool_size":    10,
            },
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tmpDir := t.TempDir()
            
            err := RenderTemplates(tt.language, tt.answers, tmpDir)
            if err != nil {
                t.Fatalf("render failed: %v", err)
            }
            
            // Verify output
            output := filepath.Join(tmpDir, "database.go")
            if _, err := os.Stat(output); os.IsNotExist(err) {
                t.Error("expected database.go to be generated")
            }
        })
    }
}
```

### Integration Tests

```bash
#!/bin/bash
# tests/integration.sh

set -e

# Build crux with plugin
crux plugin install ../

# Generate test service
crux new test-service \
    --language go \
    --framework gin \
    --db postgresql \
    --db-read-replica true

# Verify generated service
cd test-service

# Should compile
go build ./...

# Should pass tests
go test ./...

# Should have database files
test -f internal/database/database.go
test -f internal/database/repository.go

echo "Integration tests passed"
```

## Publishing Plugins

### Tier 2 Plugin Submission

1. Create plugin repository
2. Implement plugin following this guide
3. Add comprehensive tests
4. Write README.md with usage examples
5. Submit PR to plugin registry

**PR Template:**
```markdown
## Plugin Information

- Name: crux-plugin-postgresql
- Version: 2.1.0
- Tier: 2
- Author: developer@company.com

## Description

PostgreSQL integration with connection pooling, read replicas, and migration support.

## Testing

- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Tested with Go, Python, Java
- [ ] Documentation complete

## Checklist

- [ ] plugin.yaml valid
- [ ] questions.yaml valid
- [ ] Templates for all supported languages
- [ ] Validators defined
- [ ] Hooks tested
- [ ] README.md complete
- [ ] LICENSE file included
```

### Plugin Registry Entry

```yaml
# registry/plugins/postgresql.yaml
name: crux-plugin-postgresql
tier: 2
version: 2.1.0
description: PostgreSQL integration with connection pooling and migrations
author: developer@company.com
repository: https://github.com/company/crux-plugin-postgresql
homepage: https://github.com/company/crux-plugin-postgresql/blob/main/README.md
license: MIT
languages:
  - go
  - python
  - java
  - node
tags:
  - database
  - postgresql
  - sql
downloads: 1234
stars: 56
last_updated: 2026-03-09
```

## Best Practices

### Do's

- Keep plugins focused on single responsibility
- Support multiple languages when possible
- Provide sensible defaults
- Include comprehensive help text
- Write thorough tests
- Document all configuration options
- Use semantic versioning
- Handle errors gracefully

### Don'ts

- Don't modify core crux files
- Don't assume specific directory structure
- Don't hardcode paths or URLs
- Don't make network calls without user consent
- Don't include secrets in templates
- Don't use global state
- Don't break backward compatibility in minor versions

## Examples

### Minimal Plugin

```yaml
# plugin.yaml
name: crux-plugin-redis
version: 1.0.0
tier: 2
description: Redis caching integration
author: dev@company.com
crux_core_compatibility: ">=1.0.0"
languages_supported: [go, python, java, node]
capabilities: [questions, templates]
questions_schema: ./questions.yaml
templates_dir: ./templates/
```

```yaml
# questions.yaml
questions:
  - id: redis_host
    prompt: "Redis host?"
    type: text
    default: "localhost"
    group: "Cache"
    order: 10
```

```go
// templates/go/redis.go.tmpl
package cache

import "github.com/go-redis/redis/v8"

func NewRedisClient() *redis.Client {
    return redis.NewClient(&redis.Options{
        Addr: "{{ .redis_host }}:6379",
    })
}
```

## References

- Template Syntax: https://pkg.go.dev/text/template
- YAML Specification: https://yaml.org/spec/1.2/spec.html
- Semantic Versioning: https://semver.org/
