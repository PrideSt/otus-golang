package validator

import "fmt"

type ValidationError struct {
	Field string
	Err   error
}

type Validated interface {
	Validate() ([]ValidationError, error)
}

// NewError created new ValidationError for specific field and err.
func NewError(field string, err error) ValidationError {
	return ValidationError{
		Field: field,
		Err:   err,
	}
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation failed for field %s: %s", e.Field, e.Err)
}
