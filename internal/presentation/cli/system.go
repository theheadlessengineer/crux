package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/theheadlessengineer/crux/internal/app/config"
)

type checkResult struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Detail string `json:"detail,omitempty"`
}

func newSystemCommand(cfg *config.GlobalConfig) *cobra.Command {
	return &cobra.Command{
		Use:   "system",
		Short: "Check system prerequisites",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			checks := runSystemChecks()

			failed := false
			for _, c := range checks {
				if c.Status == "FAIL" {
					failed = true
					break
				}
			}

			if cfg.OutputMode == "json" {
				if err := json.NewEncoder(cmd.OutOrStdout()).Encode(checks); err != nil {
					return err
				}
			} else {
				for _, c := range checks {
					detail := ""
					if c.Detail != "" {
						detail = "  (" + c.Detail + ")"
					}
					_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%-30s %s%s\n", c.Name, c.Status, detail)
				}
			}

			if failed {
				return &exitError{code: 1, msg: "one or more prerequisite checks failed"}
			}
			return nil
		},
	}
}

func runSystemChecks() []checkResult {
	return []checkResult{
		checkGoVersion(),
		checkBinary("git"),
		checkDocker(),
		checkCruxHome(),
	}
}

func checkGoVersion() checkResult {
	ver := runtime.Version() // e.g. "go1.26.1"
	return checkResult{Name: "Go version", Status: "PASS", Detail: ver}
}

func checkBinary(name string) checkResult {
	_, err := exec.LookPath(name)
	if err != nil {
		return checkResult{Name: name, Status: "FAIL", Detail: "not found in PATH"}
	}
	return checkResult{Name: name, Status: "PASS"}
}

func checkDocker() checkResult {
	if _, err := exec.LookPath("docker"); err != nil {
		return checkResult{Name: "docker", Status: "FAIL", Detail: "not found in PATH"}
	}
	out, err := exec.Command("docker", "info", "--format", "{{.ServerVersion}}").Output()
	if err != nil {
		return checkResult{Name: "docker", Status: "FAIL", Detail: "daemon not running"}
	}
	ver := strings.TrimSpace(string(out))
	return checkResult{Name: "docker", Status: "PASS", Detail: "daemon running, server " + ver}
}

func checkCruxHome() checkResult {
	home, err := os.UserHomeDir()
	if err != nil {
		return checkResult{Name: "~/.crux/ directory", Status: "FAIL", Detail: "cannot determine home directory"}
	}
	dir := home + "/.crux"
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		if mkErr := os.MkdirAll(dir, 0o700); mkErr != nil {
			return checkResult{Name: "~/.crux/ directory", Status: "FAIL", Detail: "cannot create directory"}
		}
		return checkResult{Name: "~/.crux/ directory", Status: "PASS", Detail: "created"}
	}
	if err != nil || !info.IsDir() {
		return checkResult{Name: "~/.crux/ directory", Status: "FAIL", Detail: "path exists but is not a directory"}
	}
	// Verify writable.
	tmp, err := os.CreateTemp(dir, ".crux-write-check-*")
	if err != nil {
		return checkResult{Name: "~/.crux/ directory", Status: "FAIL", Detail: "not writable"}
	}
	_ = tmp.Close()
	_ = os.Remove(tmp.Name())
	return checkResult{Name: "~/.crux/ directory", Status: "PASS"}
}
