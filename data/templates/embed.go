// Package templates exposes the embedded template filesystem.
package templates

import "embed"

// FS is the embedded filesystem containing all .tmpl files for all languages.
//
//go:embed all:go-gin all:python-fastapi all:java-spring all:node-express
var FS embed.FS
