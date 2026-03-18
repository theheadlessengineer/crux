// Package plugins exposes the embedded plugin filesystem.
package plugins

import "embed"

// FS is the embedded filesystem containing all bundled plugin directories.
//
//go:embed all:crux-plugin-auth-jwt
//go:embed all:crux-plugin-claude-code
//go:embed all:crux-plugin-github-actions
//go:embed all:crux-plugin-kafka
//go:embed all:crux-plugin-kubernetes
//go:embed all:crux-plugin-postgresql
//go:embed all:crux-plugin-prometheus
//go:embed all:crux-plugin-redis
//go:embed all:crux-plugin-terraform-aws
//go:embed all:crux-plugin-resilience
//go:embed all:crux-plugin-spiffe
//go:embed all:crux-plugin-multitenant
//go:embed all:crux-plugin-mysql
//go:embed all:crux-plugin-mongodb
//go:embed all:crux-plugin-rabbitmq
//go:embed all:crux-plugin-gitlab-ci
//go:embed all:crux-plugin-datadog
//go:embed all:crux-plugin-terraform-gcp
//go:embed all:crux-plugin-terraform-azure
//go:embed all:crux-plugin-grpc
//go:embed all:crux-plugin-claude-api
//go:embed all:crux-plugin-openai
//go:embed all:crux-plugin-github-copilot
//go:embed all:crux-plugin-cursor
var FS embed.FS
