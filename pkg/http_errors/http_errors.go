package http_errors

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/MatheusHenrique129/bemax-backend/internal/core/errs"
)

const (
	contentTypeHeaderKey       = "Content-Type"
	contentTypeApplicationJSON = "application/json"
)

// ErrorResponse defines the structure of the JSON response returned when an error occurs.
type ErrorResponse struct {
	StatusCode int      `json:"status_code"`
	Message    string   `json:"message"`
	CauseList  []string `json:"cause_list,omitempty"`
}

// ErrorHandler writes a formatted JSON error response to the provided http.ResponseWriter.
//
// It handles three main cases:
// - If the error is of type *errs.MultiError, it extracts all sub-errors and uses the specified status code.
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
	case *errs.MultiError:
		status = e.StatusCode
		for _, subErr := range e.Errors {
			response.CauseList = append(response.CauseList, subErr.Error())
		}
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
