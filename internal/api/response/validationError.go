package response

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validator = validator.New()

func ResponseWithValidationErrors(w http.ResponseWriter, statusCode int, message string, details any) {
	ResponseWithJSON(w, statusCode, map[string]any{
		"error":   message,
		"details": details,
	})
}

func FormatValidationError(err error) map[string]string {
	errors := make(map[string]string)

	ve, ok := err.(validator.ValidationErrors)
	if !ok {
		return errors
	}

	for _, f := range ve {
		field := f.Field()

		switch f.Tag() {

		case "required":
			errors[field] = "this field is required"

		case "min":
			errors[field] = fmt.Sprintf("must be at least %s characters", f.Param())

		case "max":
			errors[field] = fmt.Sprintf("must be at most %s characters", f.Param())

		case "oneof":
			errors[field] = fmt.Sprintf("must be one of: %s", strings.ReplaceAll(f.Param(), " ", ", "))

		case "gt":
			errors[field] = fmt.Sprintf("must be greater than %s", f.Param())

		case "gte":
			errors[field] = fmt.Sprintf("must be greater than or equal to %s", f.Param())

		case "datetime":
			errors[field] = "must be a valid date in YYYY-MM-DD format"

		default:
			errors[field] = "invalid value"
		}
	}

	return errors
}
