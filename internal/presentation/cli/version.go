package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/theheadlessengineer/crux/internal/app/config"
)

type versionInfo struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildTime string `json:"buildTime"`
}

func newVersionCommand(cfg *config.GlobalConfig, version, commit, buildTime string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print crux version information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			info := versionInfo{Version: version, Commit: commit, BuildTime: buildTime}
			if cfg.OutputMode == "json" {
				return json.NewEncoder(cmd.OutOrStdout()).Encode(info)
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "crux version %s (commit: %s, built: %s)\n",
				info.Version, info.Commit, info.BuildTime)
			return nil
		},
	}
}
