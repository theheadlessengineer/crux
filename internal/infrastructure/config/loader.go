// Package config provides YAML configuration file loading for crux new.
package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ServiceConfig holds service-level fields from the config file.
type ServiceConfig struct {
	Name      string `yaml:"name"`
	Language  string `yaml:"language"`
	Framework string `yaml:"framework"`
	Team      string `yaml:"team,omitempty"`
}

// Config is the parsed representation of a crux config YAML file.
type Config struct {
	Service ServiceConfig          `yaml:"service"`
	Plugins []string               `yaml:"plugins,omitempty"`
	Answers map[string]interface{} `yaml:"answers,omitempty"`
}

// requiredServiceFields lists the service sub-fields required for --no-prompt mode.
var requiredServiceFields = []string{"name", "language", "framework"}

// Load reads and parses a YAML config file from path.
// Returns an error if the file cannot be read or is not valid YAML.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path) //nolint:gosec // path is user-supplied CLI input, intentional
	if err != nil {
		return nil, fmt.Errorf("read config file %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config file %q: %w", path, err)
	}

	return &cfg, nil
}

// ValidateForNoPrompt checks that all required fields are present when running
// in --no-prompt mode. Returns a descriptive error listing missing fields.
func ValidateForNoPrompt(cfg *Config) error {
	var missing []string

	if cfg.Service.Name == "" {
		missing = append(missing, "service.name")
	}
	if cfg.Service.Language == "" {
		missing = append(missing, "service.language")
	}
	if cfg.Service.Framework == "" {
		missing = append(missing, "service.framework")
	}

	if len(missing) > 0 {
		return fmt.Errorf("--no-prompt requires the following fields in the config file: %v", missing)
	}

	return nil
}

// ErrNoConfigForNoPrompt is returned when --no-prompt is used without --config.
var ErrNoConfigForNoPrompt = errors.New(
	"--no-prompt requires --config <file>: all required answers must be provided in the config file",
)

// _ ensures requiredServiceFields is referenced (avoids unused var lint error).
var _ = requiredServiceFields
