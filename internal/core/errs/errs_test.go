package errs_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/MatheusHenrique129/bemax-backend/internal/core/errs"
	"github.com/stretchr/testify/assert"
)

func TestMultiError_Add(t *testing.T) {
	type args struct {
		errorsToAdd []error
	}

	type wants struct {
		errorCount int
	}

	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name: "should_add_valid_error",
			args: args{
				errorsToAdd: []error{errors.New("test error")},
			},
			wants: wants{
				errorCount: 1,
			},
		},
		{
			name: "should_add_multiple_errors",
			args: args{
				errorsToAdd: []error{
					errors.New("first error"),
					errors.New("second error"),
					errors.New("third error"),
				},
			},
			wants: wants{
				errorCount: 3,
			},
		},
		{
			name: "should_not_add_nil_error",
			args: args{
				errorsToAdd: []error{nil},
			},
			wants: wants{
				errorCount: 0,
			},
		},
		{
			name: "should_add_only_non_nil_errors",
			args: args{
				errorsToAdd: []error{
					errors.New("valid error"),
					nil,
					errors.New("another valid error"),
					nil,
				},
			},
			wants: wants{
				errorCount: 2,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			multiErr := &errs.MultiError{}

			// Act
			for _, err := range test.args.errorsToAdd {
				multiErr.Add(err)
			}

			// Assert
			assert.Equal(t, test.wants.errorCount, len(multiErr.Errors))
		})
	}
}

func TestMultiError_Error(t *testing.T) {
	type args struct {
		errors []error
	}

	type wants struct {
		errorMessage string
	}

	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name: "should_return_empty_string_when_no_errors",
			args: args{
				errors: []error{},
			},
			wants: wants{
				errorMessage: "",
			},
		},
		{
			name: "should_return_single_error_message",
			args: args{
				errors: []error{errors.New("single error")},
			},
			wants: wants{
				errorMessage: "single error",
			},
		},
		{
			name: "should_concatenate_multiple_errors_with_semicolon",
			args: args{
				errors: []error{
					errors.New("first error"),
					errors.New("second error"),
				},
			},
			wants: wants{
				errorMessage: "first error; second error",
			},
		},
		{
			name: "should_concatenate_three_errors_correctly",
			args: args{
				errors: []error{
					errors.New("error one"),
					errors.New("error two"),
					errors.New("error three"),
				},
			},
			wants: wants{
				errorMessage: "error one; error two; error three",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			multiErr := &errs.MultiError{Errors: test.args.errors}

			// Act
			result := multiErr.Error()

			// Assert
			assert.Equal(t, test.wants.errorMessage, result)
		})
	}
}

func TestMultiError_HasErrors(t *testing.T) {
	type args struct {
		errors []error
	}

	type wants struct {
		hasErrors bool
	}

	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name: "should_return_false_when_no_errors",
			args: args{
				errors: []error{},
			},
			wants: wants{
				hasErrors: false,
			},
		},
		{
			name: "should_return_true_when_has_single_error",
			args: args{
				errors: []error{errors.New("test error")},
			},
			wants: wants{
				hasErrors: true,
			},
		},
		{
			name: "should_return_true_when_has_multiple_errors",
			args: args{
				errors: []error{
					errors.New("first error"),
					errors.New("second error"),
				},
			},
			wants: wants{
				hasErrors: true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			multiErr := &errs.MultiError{Errors: test.args.errors}

			// Act
			result := multiErr.HasErrors()

			// Assert
			assert.Equal(t, test.wants.hasErrors, result)
		})
	}
}

func TestNewValidationError(t *testing.T) {
	// Act
	validationErr := errs.NewValidationError()

	// Assert
	assert.NotNil(t, validationErr)
	assert.Equal(t, http.StatusBadRequest, validationErr.StatusCode)
	assert.Equal(t, 0, len(validationErr.Errors))
	assert.False(t, validationErr.HasErrors())
}

func TestNewInternalError(t *testing.T) {
	// Act
	internalErr := errs.NewInternalError()

	// Assert
	assert.NotNil(t, internalErr)
	assert.Equal(t, http.StatusInternalServerError, internalErr.StatusCode)
	assert.Equal(t, 0, len(internalErr.Errors))
	assert.False(t, internalErr.HasErrors())
}

func TestNewBusinessError(t *testing.T) {
	// Act
	businessErr := errs.NewBusinessError()

	// Assert
	assert.NotNil(t, businessErr)
	assert.Equal(t, http.StatusBadRequest, businessErr.StatusCode)
	assert.Equal(t, 0, len(businessErr.Errors))
	assert.False(t, businessErr.HasErrors())
}

func TestMultiError_IntegrationFlow(t *testing.T) {
	// This test demonstrates the full workflow of MultiError usage

	// Arrange
	multiErr := errs.NewValidationError()

	// Act & Assert - Initial state
	assert.False(t, multiErr.HasErrors())
	assert.Equal(t, "", multiErr.Error())

	// Add first error
	multiErr.Add(errors.New("validation failed"))
	assert.True(t, multiErr.HasErrors())
	assert.Equal(t, "validation failed", multiErr.Error())

	// Add second error
	multiErr.Add(errors.New("field is required"))
	assert.True(t, multiErr.HasErrors())
	assert.Equal(t, "validation failed; field is required", multiErr.Error())

	// Try to add nil error (should be ignored)
	multiErr.Add(nil)
	assert.Equal(t, 2, len(multiErr.Errors))
	assert.Equal(t, "validation failed; field is required", multiErr.Error())
}
