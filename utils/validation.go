package utils

import (
	"strings"

	"personalnote.eu/simple-go-api/models"
)

// ValidateUser validates a user struct and returns validation errors
func ValidateUser(user models.User) []string {
	var validationErrors []string

	if strings.TrimSpace(user.Name) == "" {
		validationErrors = append(validationErrors, "name is required")
	}

	if user.ID <= 0 {
		validationErrors = append(validationErrors, "id must be a positive integer")
	}

	return validationErrors
}
