package secrets

import (
	"context"
	"errors"
	"fmt"
	"os"
)

// awsClient fetches secrets from AWS Secrets Manager.
// It reads AWS_REGION from the environment; credentials come from the default
// AWS credential chain (IAM role, env vars, ~/.aws/credentials).
type awsClient struct {
	region string
}

func newAWSClient() (*awsClient, error) {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		return nil, errors.New("secrets/aws: AWS_REGION is required")
	}
	return &awsClient{region: region}, nil
}

// Get fetches the secret named key from AWS Secrets Manager.
// This is a stub — replace with the AWS SDK v2 secretsmanager client in production.
func (a *awsClient) Get(_ context.Context, key string) (string, error) {
	// TODO: replace with aws.NewConfig() + secretsmanager.GetSecretValue using a.region.
	return "", fmt.Errorf("aws-secrets-manager: Get(%q) not yet implemented — wire the AWS SDK", key)
}
