// Package resilience loads resilience.yaml and exposes typed configuration
// for timeouts, retries, circuit breakers, and bulkheads.
// Invalid configuration causes a startup failure (fail-fast principle).
package resilience
