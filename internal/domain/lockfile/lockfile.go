// Package lockfile defines the domain types and writer for .skeleton.json and crux.lock.
package lockfile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// SkeletonService holds service-level metadata recorded in .skeleton.json.
type SkeletonService struct {
	Name      string `json:"name"`
	Language  string `json:"language"`
	Framework string `json:"framework"`
}

// PluginEntry records a plugin name and its resolved version.
type PluginEntry struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// Tier1Standards records whether Tier 1 standards were enforced.
type Tier1Standards struct {
	Enforced          bool     `json:"enforced"`
	DisabledStandards []string `json:"disabledStandards"`
}

// Skeleton is the full .skeleton.json document.
type Skeleton struct {
	CruxVersion    string          `json:"cruxVersion"`
	GeneratedAt    time.Time       `json:"generatedAt"`
	Service        SkeletonService `json:"service"`
	Answers        map[string]any  `json:"answers"`
	Plugins        []PluginEntry   `json:"plugins"`
	Deviations     []string        `json:"deviations"`
	Tier1Standards Tier1Standards  `json:"tier1Standards"`
}

// Lockfile is the crux.lock document — exact plugin version pins only.
type Lockfile struct {
	LockfileVersion int               `json:"lockfileVersion"`
	Plugins         map[string]string `json:"plugins"`
}

// Write writes .skeleton.json and crux.lock into outputDir.
func Write(outputDir string, skeleton *Skeleton) error {
	if err := writeJSON(filepath.Join(outputDir, ".skeleton.json"), skeleton); err != nil {
		return fmt.Errorf("write .skeleton.json: %w", err)
	}

	lock := buildLockfile(skeleton.Plugins)
	if err := writeJSON(filepath.Join(outputDir, "crux.lock"), lock); err != nil {
		return fmt.Errorf("write crux.lock: %w", err)
	}

	return nil
}

func buildLockfile(plugins []PluginEntry) *Lockfile {
	pins := make(map[string]string, len(plugins))
	for _, p := range plugins {
		pins[p.Name] = p.Version
	}
	return &Lockfile{LockfileVersion: 1, Plugins: pins}
}

func writeJSON(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}
