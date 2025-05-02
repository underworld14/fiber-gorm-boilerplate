package validators

import (
	"github.com/go-playground/validator/v10"
)

// Validator instance
var validate = validator.New()

// Validate validates a struct based on the validator tags
func Validate(s interface{}) error {
	return validate.Struct(s)
}

// ValidationErrors returns a map of field errors
func ValidationErrors(err error) map[string]string {
	if err == nil {
		return nil
	}

	errors := make(map[string]string)
	
	for _, err := range err.(validator.ValidationErrors) {
		errors[err.Field()] = getErrorMsg(err)
	}
	
	return errors
}

// getErrorMsg returns a more user-friendly error message based on the validation tag
func getErrorMsg(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Should be at least " + err.Param() + " characters long"
	default:
		return "Invalid value"
	}
}
