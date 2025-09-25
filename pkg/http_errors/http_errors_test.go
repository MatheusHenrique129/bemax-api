package http_errors_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/MatheusHenrique129/bemax-backend/internal/core/errs"
	"github.com/MatheusHenrique129/bemax-backend/pkg/http_errors"
	"github.com/stretchr/testify/assert"
)

func TestErrorHandler(t *testing.T) {
	tests := []struct {
		name         string
		inputError   any
		expectedCode int
	}{
		{
			name: "MultiError case",
			inputError: &errs.MultiError{
				Errors:     []error{errors.New("error 1"), errors.New("error 2")},
				StatusCode: http.StatusBadRequest,
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "json.UnmarshalTypeError case",
			inputError: &json.UnmarshalTypeError{
				Value: "string",
				Type:  reflect.TypeOf(""),
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "generic error with invalid character",
			inputError:   errors.New("invalid character 'x' in JSON"),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "generic error",
			inputError:   errors.New("some other error"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			http_errors.ErrorHandler(w, tt.inputError.(error))
			assert.Equal(t, tt.expectedCode, w.Code)

		})
	}
}
