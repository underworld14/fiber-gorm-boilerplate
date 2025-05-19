package validators

import (
	"errors"
	"fiber-gorm/internal/models"
	"regexp"
)

func ValidateLogin(payload *models.LoginUserPayload) error {
	if err := Validate(payload); err != nil {
		return err
	}

	if err := ValidatePassword(payload.Password); err != nil {
		return err
	}

	return nil
}

// validatePassword checks if password meets security requirements
func ValidatePassword(password string) error {
	// Check for at least one uppercase letter
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return errors.New("password must contain at least one uppercase letter")
	}

	// Check for at least one digit
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return errors.New("password must contain at least one digit")
	}

	// Check for at least one special character
	if !regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password) {
		return errors.New("password must contain at least one special character")
	}

	return nil
}
