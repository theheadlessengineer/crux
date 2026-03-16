package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/theheadlessengineer/crux/internal/infrastructure/config"
)

const validYAML = `
service:
  name: payment-service
  language: go
  framework: gin
  team: payments-team

plugins:
  - crux-plugin-postgresql
  - crux-plugin-redis

answers:
  pg_version: "16"
  redis_mode: cluster
`

func writeFile(t *testing.T, content string) string {
	t.Helper()
	f := filepath.Join(t.TempDir(), "crux.config.yaml")
	require.NoError(t, os.WriteFile(f, []byte(content), 0o644))
	return f
}

func TestLoad_ParsesValidYAML(t *testing.T) {
	path := writeFile(t, validYAML)
	cfg, err := config.Load(path)
	require.NoError(t, err)

	assert.Equal(t, "payment-service", cfg.Service.Name)
	assert.Equal(t, "go", cfg.Service.Language)
	assert.Equal(t, "gin", cfg.Service.Framework)
	assert.Equal(t, "payments-team", cfg.Service.Team)
	assert.Equal(t, []string{"crux-plugin-postgresql", "crux-plugin-redis"}, cfg.Plugins)
	assert.Equal(t, "16", cfg.Answers["pg_version"])
	assert.Equal(t, "cluster", cfg.Answers["redis_mode"])
}

func TestLoad_ErrorOnMissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/crux.config.yaml")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "read config file")
}

func TestLoad_ErrorOnInvalidYAML(t *testing.T) {
	path := writeFile(t, "service: [invalid: yaml: :")
	_, err := config.Load(path)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "parse config file")
}

func TestLoad_EmptyAnswersAndPlugins(t *testing.T) {
	path := writeFile(t, `
service:
  name: my-service
  language: go
  framework: gin
`)
	cfg, err := config.Load(path)
	require.NoError(t, err)
	assert.Empty(t, cfg.Plugins)
	assert.Empty(t, cfg.Answers)
}

func TestValidateForNoPrompt_PassesWhenAllFieldsPresent(t *testing.T) {
	cfg := &config.Config{
		Service: config.ServiceConfig{Name: "svc", Language: "go", Framework: "gin"},
	}
	assert.NoError(t, config.ValidateForNoPrompt(cfg))
}

func TestValidateForNoPrompt_ErrorOnMissingName(t *testing.T) {
	cfg := &config.Config{
		Service: config.ServiceConfig{Language: "go", Framework: "gin"},
	}
	err := config.ValidateForNoPrompt(cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "service.name")
}

func TestValidateForNoPrompt_ErrorOnMissingLanguage(t *testing.T) {
	cfg := &config.Config{
		Service: config.ServiceConfig{Name: "svc", Framework: "gin"},
	}
	err := config.ValidateForNoPrompt(cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "service.language")
}

func TestValidateForNoPrompt_ErrorOnMissingFramework(t *testing.T) {
	cfg := &config.Config{
		Service: config.ServiceConfig{Name: "svc", Language: "go"},
	}
	err := config.ValidateForNoPrompt(cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "service.framework")
}

func TestValidateForNoPrompt_ErrorListsAllMissingFields(t *testing.T) {
	cfg := &config.Config{}
	err := config.ValidateForNoPrompt(cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "service.name")
	assert.Contains(t, err.Error(), "service.language")
	assert.Contains(t, err.Error(), "service.framework")
}
