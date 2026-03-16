// Package templates exposes the embedded template filesystem.
package templates

import "embed"

// FS is the embedded filesystem containing all .tmpl files.
//
//go:embed all:go-gin
var FS embed.FS
