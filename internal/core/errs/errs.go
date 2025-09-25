package errs

import (
	"net/http"
	"strings"
)

const validationError = http.StatusBadRequest        // HTTP status code for validation errors
const internalError = http.StatusInternalServerError // HTTP status code for internal server errors
const businessError = http.StatusBadRequest          // HTTP status code for business logic errors

// MultiError aggregates multiple errors and associates them with a status code.
// It implements the standard Go error interface.
type MultiError struct {
	Errors     []error
	StatusCode int
}

// Add appends a new error to the MultiError collection if it is not nil.
func (m *MultiError) Add(err error) {
	if err != nil {
		m.Errors = append(m.Errors, err)
	}
}

// Error returns a single concatenated string of all errors in the collection.
// Each error message is separated by '; '.
func (m *MultiError) Error() string {
	if len(m.Errors) == 0 {
		return ""
	}
	var sb strings.Builder
	for _, err := range m.Errors {
		sb.WriteString(err.Error())
		sb.WriteString("; ")
	}
	return strings.TrimSuffix(sb.String(), "; ")
}

// HasErrors checks whether the MultiError contains any errors.
func (m *MultiError) HasErrors() bool {
	return len(m.Errors) > 0
}

// NewValidationError creates a new MultiError instance categorized as a validation error (HTTP 400).
func NewValidationError() *MultiError { return &MultiError{StatusCode: validationError} }

// NewInternalError creates a new MultiError instance categorized as an internal server error (HTTP 500).
func NewInternalError() *MultiError {
	return &MultiError{StatusCode: internalError}
}

// NewBusinessError creates a new MultiError instance categorized as a business logic error (HTTP 400).
func NewBusinessError() *MultiError {
	return &MultiError{StatusCode: businessError}
}
