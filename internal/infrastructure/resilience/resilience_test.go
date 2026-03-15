package resilience

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeYAML(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "resilience-*.yaml")
	require.NoError(t, err)
	_, err = f.WriteString(content)
	require.NoError(t, err)
	require.NoError(t, f.Close())
	return f.Name()
}

func TestLoad_ValidFile_ReturnsConfig(t *testing.T) {
	path := writeYAML(t, `
timeout:
  defaultMs: 3000
retry:
  maxAttempts: 5
  backoffType: exponential
  initialIntervalMs: 200
  maxIntervalMs: 2000
circuitBreaker:
  failureRateThreshold: 60
  slowCallRateThreshold: 80
  minimumCallCount: 5
  openStateWaitDurationMs: 30000
bulkhead:
  maxConcurrentCalls: 10
`)

	cfg, err := Load(path)

	require.NoError(t, err)
	assert.Equal(t, 3000*time.Millisecond, cfg.TimeoutDuration())
	assert.Equal(t, 5, cfg.Retry.MaxAttempts)
	assert.Equal(t, "exponential", cfg.Retry.BackoffType)
	assert.Equal(t, 200*time.Millisecond, cfg.InitialBackoff())
	assert.Equal(t, 2000*time.Millisecond, cfg.MaxBackoff())
	assert.InDelta(t, 60.0, cfg.CircuitBreaker.FailureRateThreshold, 0.001)
	assert.Equal(t, 30000*time.Millisecond, cfg.OpenStateWait())
	assert.Equal(t, 10, cfg.Bulkhead.MaxConcurrentCalls)
}

func TestLoad_MissingFile_ReturnsError(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "nonexistent.yaml"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "resilience:")
}

func TestLoad_InvalidYAML_ReturnsError(t *testing.T) {
	path := writeYAML(t, "not: valid: yaml: [")
	_, err := Load(path)
	require.Error(t, err)
}

func TestLoad_InvalidTimeout_ReturnsError(t *testing.T) {
	path := writeYAML(t, `
timeout:
  defaultMs: 0
retry:
  maxAttempts: 3
  initialIntervalMs: 100
circuitBreaker:
  failureRateThreshold: 50
bulkhead:
  maxConcurrentCalls: 10
`)
	_, err := Load(path)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "timeout.defaultMs")
}

func TestLoad_InvalidCircuitBreakerThreshold_ReturnsError(t *testing.T) {
	path := writeYAML(t, `
timeout:
  defaultMs: 5000
retry:
  maxAttempts: 3
  initialIntervalMs: 100
circuitBreaker:
  failureRateThreshold: 150
bulkhead:
  maxConcurrentCalls: 10
`)
	_, err := Load(path)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failureRateThreshold")
}

func TestLoad_EmptyFile_UsesDefaults(t *testing.T) {
	path := writeYAML(t, "")

	cfg, err := Load(path)

	require.NoError(t, err)
	assert.Equal(t, 5000*time.Millisecond, cfg.TimeoutDuration())
	assert.Equal(t, 3, cfg.Retry.MaxAttempts)
	assert.Equal(t, 20, cfg.Bulkhead.MaxConcurrentCalls)
}
