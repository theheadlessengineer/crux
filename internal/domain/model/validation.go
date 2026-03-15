// Package model provides core domain models and validation.
package model

import (
	"fmt"
	"regexp"
)

// ValidateServiceName validates a service name according to Kubernetes naming conventions.
func ValidateServiceName(name string) error {
	if name == "" {
		return fmt.Errorf("service name cannot be empty")
	}

	// Pattern: ^[a-z][a-z0-9-]{2,62}$ — minimum 3 chars total, maximum 63.
	matched, _ := regexp.MatchString(`^[a-z][a-z0-9-]{2,62}$`, name)
	if !matched {
		return fmt.Errorf(
			"invalid service name format: must start with a lowercase letter, " +
				"be at least 3 characters, and contain only lowercase letters, numbers, and hyphens",
		)
	}

	return nil
}
