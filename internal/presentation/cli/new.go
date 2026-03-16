package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/theheadlessengineer/crux/internal/app/config"
	"github.com/theheadlessengineer/crux/internal/domain/lockfile"
	"github.com/theheadlessengineer/crux/internal/domain/model"
	infraconfig "github.com/theheadlessengineer/crux/internal/infrastructure/config"
	"github.com/theheadlessengineer/crux/internal/infrastructure/generator"
)

type newCommandFlags struct {
	outputDir  string
	configFile string
	dryRun     bool
	noPrompt   bool
}

func newNewCommand(_ *config.GlobalConfig, cruxVersion string) *cobra.Command {
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

			var fileCfg *infraconfig.Config
			if flags.configFile != "" {
				loaded, err := infraconfig.Load(flags.configFile)
				if err != nil {
					return fmt.Errorf("load config: %w", err)
				}
				fileCfg = loaded
			}

			if flags.noPrompt {
				if fileCfg == nil {
					return infraconfig.ErrNoConfigForNoPrompt
				}
				if err := infraconfig.ValidateForNoPrompt(fileCfg); err != nil {
					return err
				}
			}

			if flags.dryRun {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "[dry-run] Would generate service: %s\n", name)
				return nil
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Generating %s...\n", name)

			outDir := flags.outputDir
			if outDir == "" {
				outDir = name
			}

			if err := os.MkdirAll(outDir, 0o750); err != nil {
				return fmt.Errorf("create output directory %q: %w", outDir, err)
			}

			genCfg := buildGeneratorConfig(name, cruxVersion, fileCfg)
			if err := generator.Generate(cmd.Context(), &genCfg, outDir); err != nil {
				return fmt.Errorf("generate skeleton: %w", err)
			}

			skel := buildSkeleton(name, cruxVersion, fileCfg)
			if err := lockfile.Write(outDir, skel); err != nil {
				return fmt.Errorf("write lockfiles: %w", err)
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "✔  skeleton generated\n")
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "✔  .skeleton.json written\n")
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "✔  crux.lock written\n")
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

func buildGeneratorConfig(name, cruxVersion string, fileCfg *infraconfig.Config) generator.Config {
	cfg := generator.Config{
		ServiceName: name,
		Language:    "go",
		Framework:   "gin",
		CLIVersion:  cruxVersion,
		GeneratedAt: time.Now().UTC(),
	}
	if fileCfg != nil {
		if fileCfg.Service.Language != "" {
			cfg.Language = fileCfg.Service.Language
		}
		if fileCfg.Service.Framework != "" {
			cfg.Framework = fileCfg.Service.Framework
		}
		if fileCfg.Service.Team != "" {
			cfg.Team = fileCfg.Service.Team
		}
	}
	return cfg
}

func buildSkeleton(name, cruxVersion string, fileCfg *infraconfig.Config) *lockfile.Skeleton {
	skel := &lockfile.Skeleton{
		CruxVersion: cruxVersion,
		GeneratedAt: time.Now().UTC(),
		Service: lockfile.SkeletonService{
			Name: name,
		},
		Answers:    make(map[string]any),
		Plugins:    []lockfile.PluginEntry{},
		Deviations: []string{},
		Tier1Standards: lockfile.Tier1Standards{
			Enforced:          true,
			DisabledStandards: []string{},
		},
	}

	if fileCfg != nil {
		if fileCfg.Service.Language != "" {
			skel.Service.Language = fileCfg.Service.Language
		}
		if fileCfg.Service.Framework != "" {
			skel.Service.Framework = fileCfg.Service.Framework
		}
		for k, v := range fileCfg.Answers {
			skel.Answers[k] = v
		}
	}

	return skel
}
