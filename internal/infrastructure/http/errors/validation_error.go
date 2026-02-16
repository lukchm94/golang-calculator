package reqErr

import "fmt"

// MissingFieldError is a struct because it carries data (FieldName)
type MissingFieldError struct {
	FieldName string
}

func (e MissingFieldError) Error() string {
	return fmt.Sprintf("missing required field: %s", e.FieldName)
}

// InvalidRequestError for general structural issues
type InvalidRequestError struct {
	Details string
}

func (e InvalidRequestError) Error() string {
	return fmt.Sprintf("invalid request: %s", e.Details)
}

type InvalidRequestMethodError struct {
	Method string
}

func (e InvalidRequestMethodError) Error() string {
	return fmt.Sprintf("invalid HTTP method: %s", e.Method)
}
