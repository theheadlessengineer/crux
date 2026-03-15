package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/theheadlessengineer/crux/internal/app/config"
)

// tier1RequiredFiles are the Tier 1 files every generated service must contain.
var tier1RequiredFiles = []string{
	"README.md",
	"Makefile",
	".skeleton.json",
}

type fileCheck struct {
	File    string `json:"file"`
	Present bool   `json:"present"`
}

func newValidateCommand(cfg *config.GlobalConfig) *cobra.Command {
	return &cobra.Command{
		Use:   "validate [directory]",
		Short: "Validate a generated service for Tier 1 compliance",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := "."
			if len(args) == 1 {
				dir = args[0]
			}

			results := make([]fileCheck, 0, len(tier1RequiredFiles))
			missing := 0
			for _, f := range tier1RequiredFiles {
				_, err := os.Stat(filepath.Join(dir, f))
				present := err == nil
				if !present {
					missing++
				}
				results = append(results, fileCheck{File: f, Present: present})
			}

			if cfg.OutputMode == "json" {
				return json.NewEncoder(cmd.OutOrStdout()).Encode(results)
			}

			for _, r := range results {
				mark := "✔"
				if !r.Present {
					mark = "✘"
				}
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s  %s\n", mark, r.File)
			}

			if missing > 0 {
				return &exitError{code: 1, msg: fmt.Sprintf("%d required Tier 1 file(s) missing", missing)}
			}
			return nil
		},
	}
}
