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
var FS embed.FS
