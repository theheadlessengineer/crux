// Package secrets provides a startup secrets loader that fetches secrets from
// Vault (AppRole / Kubernetes auth) or AWS Secrets Manager.
//
// The backend is selected via the SECRETS_BACKEND environment variable
// ("vault" or "aws-secrets-manager"). The service fails fast on startup if
// required secrets cannot be fetched.
package secrets
