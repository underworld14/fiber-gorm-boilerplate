package validators

import (
	"errors"
	"fiber-gorm/internal/models"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ValidateUserCreation performs custom validations beyond the struct tags
func ValidateUserCreation(payload *models.CreateUserPayload) error {
	// Basic validation using struct tags
	if err := Validate(payload); err != nil {
		return err
	}

	// Custom validations
	if err := ValidatePassword(payload.Password); err != nil {
		return err
	}

	if err := validateName(payload.Name); err != nil {
		return err
	}

	return nil
}

// validateName checks if name is valid
func validateName(name string) error {
	name = strings.TrimSpace(name)

	if len(name) < 2 {
		return errors.New("name must be at least 2 characters long")
	}

	// Check if name contains only letters and spaces
	if !regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString(name) {
		return errors.New("name must contain only letters and spaces")
	}

	return nil
}

// FormatValidationError formats a validation error into a map for API response
func FormatValidationError(err error, obj interface{}) interface{} {
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		return ValidationErrors(validationErrs, obj)
	}

	// Handle custom validation errors
	return err.Error()
}
