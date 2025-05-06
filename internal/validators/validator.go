package validators

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator instance
var validate = validator.New()

// Validate validates a struct based on the validator tags
func Validate(s interface{}) error {
	return validate.Struct(s)
}

// ValidationErrors returns a map of field errors
func ValidationErrors(err error, obj interface{}) map[string]string {
	if err == nil {
		return nil
	}

	errors := make(map[string]string)

	jsonFieldMap := map[string]string{}
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" && jsonTag != "-" {
			jsonFieldMap[field.Name] = jsonTag
		}
	}

	for _, err := range err.(validator.ValidationErrors) {
		fieldName := err.Field()
		jsonField := jsonFieldMap[fieldName]
		if jsonField != "" {
			jsonField = strings.ToLower(jsonField) // fallback to field name if json tag is not found
		}
		errors[jsonField] = getErrorMsg(err)
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
