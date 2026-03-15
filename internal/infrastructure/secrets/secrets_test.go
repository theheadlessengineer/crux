package secrets

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockClient is a test double for Client.
type mockClient struct {
	data map[string]string
}

func (m *mockClient) Get(_ context.Context, key string) (string, error) {
	v, ok := m.data[key]
	if !ok {
		return "", errors.New("not found")
	}
	return v, nil
}

func TestLoad_AllKeysPresent_ReturnsConfig(t *testing.T) {
	client := &mockClient{data: map[string]string{
		"db/password": "s3cr3t",
		"api/key":     "abc123",
	}}

	cfg, err := Load(context.Background(), client, []string{"db/password", "api/key"})

	require.NoError(t, err)
	assert.Equal(t, "s3cr3t", cfg.Values["db/password"])
	assert.Equal(t, "abc123", cfg.Values["api/key"])
}

func TestLoad_MissingKey_ReturnsErrMissingSecret(t *testing.T) {
	client := &mockClient{data: map[string]string{}}

	_, err := Load(context.Background(), client, []string{"db/password"})

	require.Error(t, err)
	var missing *ErrMissingSecret
	assert.True(t, errors.As(err, &missing), "error must be ErrMissingSecret")
	assert.Equal(t, "db/password", missing.Key)
}

func TestLoad_EmptyKeys_ReturnsEmptyConfig(t *testing.T) {
	client := &mockClient{data: map[string]string{}}

	cfg, err := Load(context.Background(), client, nil)

	require.NoError(t, err)
	assert.Empty(t, cfg.Values)
}

func TestNewClient_UnknownBackend_ReturnsError(t *testing.T) {
	t.Setenv(envBackend, "unknown-backend")

	_, err := NewClient()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "SECRETS_BACKEND")
}

func TestNewClient_VaultBackend_MissingAddr_ReturnsError(t *testing.T) {
	t.Setenv(envBackend, string(BackendVault))
	t.Setenv("VAULT_ADDR", "")
	t.Setenv("VAULT_TOKEN", "")

	_, err := NewClient()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "VAULT_ADDR")
}

func TestNewClient_AWSBackend_MissingRegion_ReturnsError(t *testing.T) {
	t.Setenv(envBackend, string(BackendAWSSM))
	t.Setenv("AWS_REGION", "")

	_, err := NewClient()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "AWS_REGION")
}

func TestErrMissingSecret_Message(t *testing.T) {
	err := &ErrMissingSecret{Key: "my/secret"}
	assert.Contains(t, err.Error(), "my/secret")
}
