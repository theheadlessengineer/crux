package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/theheadlessengineer/crux/internal/app/config"
	"github.com/theheadlessengineer/crux/internal/domain/model"
)

type newCommandFlags struct {
	outputDir  string
	configFile string
	dryRun     bool
	noPrompt   bool
}

func newNewCommand(cfg *config.GlobalConfig) *cobra.Command {
	flags := &newCommandFlags{}

	cmd := &cobra.Command{
		Use:   "new <service-name>",
		Short: "Generate a new microservice skeleton",
		Long:  "Initiates the service generation flow for a new microservice.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			if err := model.ValidateServiceName(name); err != nil {
				return &ValidationError{Msg: err.Error()}
			}

			if flags.dryRun {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "[dry-run] Would generate service: %s\n", name)
				return nil
			}

			// Stub: prompt engine and template engine invoked here in Epic 1.3 / 1.4.
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Generating %s...\n", name)
			_ = cfg
			_ = flags.outputDir
			_ = flags.configFile
			_ = flags.noPrompt
			return nil
		},
	}

	cmd.Flags().StringVar(&flags.outputDir, "output-dir", "",
		"Directory to write the generated service (default: ./<service-name>)")
	cmd.Flags().StringVar(&flags.configFile, "config", "", "Path to a pre-filled configuration YAML file")
	cmd.Flags().BoolVar(&flags.dryRun, "dry-run", false, "Print what would be generated without writing files")
	cmd.Flags().BoolVar(&flags.noPrompt, "no-prompt", false, "Run non-interactively using config file or defaults only")

	return cmd
}
