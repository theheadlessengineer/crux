package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateServiceName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid service name",
			input:   "payment-service",
			wantErr: false,
		},
		{
			name:    "valid single word",
			input:   "payments",
			wantErr: false,
		},
		{
			name:    "valid with numbers",
			input:   "service-v2",
			wantErr: false,
		},
		{
			name:    "empty name",
			input:   "",
			wantErr: true,
			errMsg:  "service name cannot be empty",
		},
		{
			name: "name too long",
			input: "this-is-a-very-long-service-name-that-exceeds-" +
				"the-maximum-allowed-length-of-sixty-three-characters",
			wantErr: true,
			errMsg:  "invalid service name format",
		},
		{
			name:    "starts with uppercase",
			input:   "Payment-service",
			wantErr: true,
			errMsg:  "invalid service name format",
		},
		{
			name:    "starts with number",
			input:   "1-service",
			wantErr: true,
			errMsg:  "invalid service name format",
		},
		{
			name:    "contains underscore",
			input:   "payment_service",
			wantErr: true,
			errMsg:  "invalid service name format",
		},
		{
			name:    "contains special characters",
			input:   "payment@service",
			wantErr: true,
			errMsg:  "invalid service name format",
		},
		{
			name:    "too short",
			input:   "ab",
			wantErr: true,
			errMsg:  "invalid service name format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateServiceName(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
