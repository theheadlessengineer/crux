package resilience

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config is the typed representation of resilience.yaml.
type Config struct {
	Timeout        TimeoutConfig        `yaml:"timeout"`
	Retry          RetryConfig          `yaml:"retry"`
	CircuitBreaker CircuitBreakerConfig `yaml:"circuitBreaker"`
	Bulkhead       BulkheadConfig       `yaml:"bulkhead"`
}

// TimeoutConfig holds per-integration call timeouts.
type TimeoutConfig struct {
	DefaultMS int `yaml:"defaultMs"`
}

// RetryConfig controls retry behaviour for outbound calls.
type RetryConfig struct {
	MaxAttempts       int    `yaml:"maxAttempts"`
	BackoffType       string `yaml:"backoffType"` // "exponential" or "fixed"
	InitialIntervalMS int    `yaml:"initialIntervalMs"`
	MaxIntervalMS     int    `yaml:"maxIntervalMs"`
}

// CircuitBreakerConfig controls the circuit breaker.
type CircuitBreakerConfig struct {
	FailureRateThreshold    float64 `yaml:"failureRateThreshold"`  // percentage 0–100
	SlowCallRateThreshold   float64 `yaml:"slowCallRateThreshold"` // percentage 0–100
	MinimumCallCount        int     `yaml:"minimumCallCount"`
	OpenStateWaitDurationMS int     `yaml:"openStateWaitDurationMs"`
}

// BulkheadConfig limits concurrent calls per downstream.
type BulkheadConfig struct {
	MaxConcurrentCalls int `yaml:"maxConcurrentCalls"`
}

// defaults returns a Config with sensible production defaults.
func defaults() Config {
	return Config{
		Timeout: TimeoutConfig{DefaultMS: 5000},
		Retry: RetryConfig{
			MaxAttempts:       3,
			BackoffType:       "exponential",
			InitialIntervalMS: 100,
			MaxIntervalMS:     1000,
		},
		CircuitBreaker: CircuitBreakerConfig{
			FailureRateThreshold:    50,
			SlowCallRateThreshold:   100,
			MinimumCallCount:        10,
			OpenStateWaitDurationMS: 60000,
		},
		Bulkhead: BulkheadConfig{MaxConcurrentCalls: 20},
	}
}

// Load reads resilience.yaml from path and returns a validated Config.
// Returns an error if the file is missing, unparseable, or contains invalid values.
// path is caller-controlled and must be a trusted, application-defined value.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path) //nolint:gosec // path is application-controlled, not user input
	if err != nil {
		return nil, fmt.Errorf("resilience: cannot read %s: %w", path, err)
	}

	cfg := defaults()
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("resilience: cannot parse %s: %w", path, err)
	}

	if err := validate(&cfg); err != nil {
		return nil, fmt.Errorf("resilience: invalid configuration: %w", err)
	}

	return &cfg, nil
}

// TimeoutDuration returns the default timeout as a time.Duration.
func (c *Config) TimeoutDuration() time.Duration {
	return time.Duration(c.Timeout.DefaultMS) * time.Millisecond
}

// InitialBackoff returns the initial retry backoff as a time.Duration.
func (c *Config) InitialBackoff() time.Duration {
	return time.Duration(c.Retry.InitialIntervalMS) * time.Millisecond
}

// MaxBackoff returns the maximum retry backoff as a time.Duration.
func (c *Config) MaxBackoff() time.Duration {
	return time.Duration(c.Retry.MaxIntervalMS) * time.Millisecond
}

// OpenStateWait returns the circuit breaker open-state wait as a time.Duration.
func (c *Config) OpenStateWait() time.Duration {
	return time.Duration(c.CircuitBreaker.OpenStateWaitDurationMS) * time.Millisecond
}

func validate(c *Config) error {
	if c.Timeout.DefaultMS <= 0 {
		return fmt.Errorf("timeout.defaultMs must be > 0")
	}
	if c.Retry.MaxAttempts <= 0 {
		return fmt.Errorf("retry.maxAttempts must be > 0")
	}
	if c.Retry.InitialIntervalMS <= 0 {
		return fmt.Errorf("retry.initialIntervalMs must be > 0")
	}
	if c.CircuitBreaker.FailureRateThreshold < 0 || c.CircuitBreaker.FailureRateThreshold > 100 {
		return fmt.Errorf("circuitBreaker.failureRateThreshold must be between 0 and 100")
	}
	if c.Bulkhead.MaxConcurrentCalls <= 0 {
		return fmt.Errorf("bulkhead.maxConcurrentCalls must be > 0")
	}
	return nil
}
