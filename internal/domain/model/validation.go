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

	if len(name) > 63 {
		return fmt.Errorf("service name too long (max 63 characters)")
	}

	matched, _ := regexp.MatchString(`^[a-z][a-z0-9-]*$`, name)
	if !matched {
		return fmt.Errorf(
			"invalid service name format: must start with lowercase letter " +
				"and contain only lowercase letters, numbers, and hyphens",
		)
	}

	return nil
}
