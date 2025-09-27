package apierrors

import (
	"fmt"
	"net/http"
)

type CauseList []interface{}

type RestError interface {
	Message() string
	Error() string
	Code() string
	Status() int
	Cause() CauseList
}

type restError struct {
	ErrorMessage string    `json:"message"`
	ErrorCode    string    `json:"error"`
	ErrorCause   CauseList `json:"causes,omitempty"`
	ErrorStatus  int       `json:"status"`
}

func (r restError) Message() string {
	return r.ErrorMessage
}

func (r restError) Code() string {
	return r.ErrorCode
}

func (r restError) Status() int {
	return r.ErrorStatus
}

func (r restError) Cause() CauseList {
	return r.ErrorCause
}

func (c CauseList) ToString() string {
	return fmt.Sprint(c)
}

func (r restError) Error() string {
	return fmt.Sprintf("Message: %s;Error Code: %s;Status: %d;CauseList: %v", r.ErrorMessage, r.ErrorCode, r.ErrorStatus, r.ErrorCause)
}

func NewRestError(message string, error string, status int, cause CauseList) RestError {
	return restError{message, error, cause, status}
}

func NewNotFoundRestError(message string) RestError {
	return restError{message, "not_found", nil, http.StatusNotFound}
}

func NewBadRequestRestError(message string) RestError {
	return restError{message, "bad_request", nil, http.StatusBadRequest}
}

func NewBadRequestValidationRestError(message string, error string, cause CauseList) RestError {
	return restError{message, error, cause, http.StatusBadRequest}
}

func NewUnauthorizedRestError(message string) RestError {
	return restError{message, "unauthorized", nil, http.StatusUnauthorized}
}

func NewTooManyRequestsRestError(message string) RestError {
	return restError{message, "too_many_requests", nil, http.StatusTooManyRequests}
}

func NewForbiddenRestError(message string) RestError {
	return restError{message, "forbidden", nil, http.StatusForbidden}
}

func NewConflictRestError(message string, cause CauseList) RestError {
	return restError{message, "conflict_error", cause, http.StatusConflict}
}

func NewInternalServerRestError(message string, err error) RestError {
	cause := CauseList{}
	if err != nil {
		cause = append(cause, err.Error())
	}
	return restError{message, "internal_server_error", cause, http.StatusInternalServerError}
}
