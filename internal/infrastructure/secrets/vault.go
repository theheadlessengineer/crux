package secrets

import (
	"context"
	"errors"
	"fmt"
	"os"
)

// vaultClient fetches secrets from HashiCorp Vault.
// It reads VAULT_ADDR and VAULT_TOKEN from the environment.
// For production use, replace token auth with AppRole or Kubernetes auth.
type vaultClient struct {
	addr  string
	token string
}

func newVaultClient() (*vaultClient, error) {
	addr := os.Getenv("VAULT_ADDR")
	if addr == "" {
		return nil, errors.New("secrets/vault: VAULT_ADDR is required")
	}
	token := os.Getenv("VAULT_TOKEN")
	if token == "" {
		return nil, errors.New("secrets/vault: VAULT_TOKEN is required")
	}
	return &vaultClient{addr: addr, token: token}, nil
}

// Get fetches the secret at path key from Vault KV v2.
// The key is expected to be in the form "<mount>/<path>#<field>".
// This is a stub — replace with the official Vault SDK in production.
func (v *vaultClient) Get(_ context.Context, key string) (string, error) {
	// TODO: replace with vault.NewClient() + KV v2 read using v.addr / v.token.
	return "", fmt.Errorf("vault: Get(%q) not yet implemented — wire the Vault SDK", key)
}
