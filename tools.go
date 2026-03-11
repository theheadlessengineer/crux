//go:build tools
// +build tools

// Package tools tracks tool dependencies for the project.
// This file ensures dependencies are recorded in go.mod even if not directly imported.
package tools

import (
	_ "github.com/spf13/cobra"
	_ "github.com/stretchr/testify/assert"
)
