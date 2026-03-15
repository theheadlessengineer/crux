// Package main is the entrypoint for the crux CLI tool.
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/theheadlessengineer/crux/internal/presentation/cli"
)

// Injected via ldflags: -X main.version=x.y.z -X main.commit=abc -X main.buildTime=...
var (
	version   = "dev"
	commit    = "none"
	buildTime = "unknown"
)

func main() {
	root := cli.BuildRoot(version, commit, buildTime)

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)

		var ve *cli.ValidationError
		if errors.As(err, &ve) {
			os.Exit(2)
		}
		os.Exit(1)
	}
}
