package http_errors

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/MatheusHenrique129/bemax-api/internal/core/apierrors"
)

const (
	contentTypeHeaderKey       = "Content-Type"
	contentTypeApplicationJSON = "application/json"
)

// ErrorResponse defines the structure of the JSON response returned when an error occurs.
type ErrorResponse struct {
	Message    string              `json:"message"`
	ErrorCode  string              `json:"error"`
	CauseList  apierrors.CauseList `json:"causes,omitempty"`
	StatusCode int                 `json:"status_code"`
}

// ErrorHandler writes a formatted JSON error response to the provided http.ResponseWriter.
//
// It handles three main cases:
// - If the error is of type *apierrors.RestError, it extracts all sub-errors and uses the specified status code.
// - If the error message indicates a JSON parsing error, it returns a 400 Bad Request with a cause "invalid json format".
// - For all other errors, it returns a 500 Internal Server Error with a cause "internal error".
//
// Parameters:
// - w: the HTTP response writer to write the error response to.
// - err: the error to be handled and formatted.
//
// Returns:
// - Always returns nil (as it writes the response directly).
func ErrorHandler(w http.ResponseWriter, err error) {
	w.Header().Set(contentTypeHeaderKey, contentTypeApplicationJSON)

	response := ErrorResponse{
		Message: "An error occurred",
	}

	status := http.StatusInternalServerError

	switch e := err.(type) {
	case apierrors.RestError:
		response.Message = e.Message()
		response.ErrorCode = e.Code()
		status = e.Status()
		response.CauseList = e.Cause()
	case *json.UnmarshalTypeError:
		status = http.StatusBadRequest
		response.CauseList = append(response.CauseList, "invalid json format")
	default:
		errDescription := err.Error()
		if strings.Contains(errDescription, "invalid parameter") ||
			strings.Contains(errDescription, "invalid character") {
			status = http.StatusBadRequest
			response.CauseList = append(response.CauseList, "invalid json format")
		} else {
			response.CauseList = append(response.CauseList, "internal error")
		}
	}

	response.StatusCode = status
	w.WriteHeader(status)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), status)
	}
}
