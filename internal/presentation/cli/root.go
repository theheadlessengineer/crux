// Package cli wires the Cobra root command and all subcommands.
package cli

import (
	"github.com/spf13/cobra"
	"github.com/theheadlessengineer/crux/internal/app/config"
)

// BuildRoot constructs and returns the root Cobra command.
// version, commit, and buildTime are injected at build time via ldflags.
func BuildRoot(version, commit, buildTime string) *cobra.Command {
	cfg := &config.GlobalConfig{}

	root := &cobra.Command{
		Use:     "crux",
		Short:   "crux — Microservice skeleton generator",
		Long:    "crux generates production-ready, company-compliant microservice skeletons.",
		Version: version,
		// Silence default error printing — we handle it ourselves.
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	root.PersistentFlags().BoolVarP(&cfg.Verbose, "verbose", "v", false, "Enable verbose output")
	root.PersistentFlags().StringVar(&cfg.OutputMode, "output", "text", "Output format: text or json")
	root.PersistentFlags().StringVar(&cfg.ConfigFile, "config", "", "Path to a pre-filled configuration YAML file")

	root.AddCommand(newNewCommand(cfg))
	root.AddCommand(newVersionCommand(cfg, version, commit, buildTime))
	root.AddCommand(newSystemCommand(cfg))
	root.AddCommand(newValidateCommand(cfg))

	return root
}
