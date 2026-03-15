package secrets

import (
	"context"
	"errors"
	"fmt"
	"os"
)

// Backend identifies the secrets provider.
type Backend string

const (
	BackendVault    Backend = "vault"
	BackendAWSSM    Backend = "aws-secrets-manager"
	envBackend              = "SECRETS_BACKEND"
)

// ErrMissingSecret is returned when a required secret cannot be found.
type ErrMissingSecret struct {
	Key string
}

func (e *ErrMissingSecret) Error() string {
	return fmt.Sprintf("required secret %q is missing", e.Key)
}

// Client is the interface that both Vault and AWS SM adapters implement.
type Client interface {
	// Get returns the value for key, or an error if the secret does not exist.
	Get(ctx context.Context, key string) (string, error)
}

// Config holds the secrets fetched at startup.
// Add fields here for each secret the service requires.
type Config struct {
	// Values is a map of secret key → value populated by Load.
	Values map[string]string
}

// Load fetches all required keys from the configured backend.
// It returns ErrMissingSecret (wrapped) for the first missing key.
// Secret values are never logged.
func Load(ctx context.Context, client Client, requiredKeys []string) (*Config, error) {
	cfg := &Config{Values: make(map[string]string, len(requiredKeys))}
	for _, key := range requiredKeys {
		val, err := client.Get(ctx, key)
		if err != nil {
			return nil, fmt.Errorf("secrets.Load: %w", &ErrMissingSecret{Key: key})
		}
		cfg.Values[key] = val
	}
	return cfg, nil
}

// NewClient returns the Client for the backend declared in SECRETS_BACKEND.
// Returns an error if the backend is unknown or required env vars are absent.
func NewClient() (Client, error) {
	switch Backend(os.Getenv(envBackend)) {
	case BackendVault:
		return newVaultClient()
	case BackendAWSSM:
		return newAWSClient()
	default:
		return nil, errors.New("secrets: SECRETS_BACKEND must be \"vault\" or \"aws-secrets-manager\"")
	}
}
