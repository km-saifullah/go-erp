package errors

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code       string
	Message    string
	HTTPStatus int
	Err        error
	Details    any
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}

	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

const (
	CodeBadRequest = "BAD_REQUEST"

	CodeValidation = "VALIDATION_ERROR"

	CodeUnauthorized = "UNAUTHORIZED"
	CodeForbidden    = "FORBIDDEN"

	CodeNotFound = "NOT_FOUND"

	CodeConflict = "CONFLICT"

	CodeTooManyRequests = "TOO_MANY_REQUESTS"

	CodeInternalServer = "INTERNAL_SERVER_ERROR"

	CodeServiceUnavailable = "SERVICE_UNAVAILABLE"
)

func New(
	code string,
	message string,
	status int,
) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: status,
	}
}

func BadRequest(message string) *AppError {
	return New(
		CodeBadRequest,
		message,
		http.StatusBadRequest,
	)
}

func Validation(message string, details any) *AppError {
	return &AppError{
		Code:       CodeValidation,
		Message:    message,
		HTTPStatus: http.StatusBadRequest,
		Details:    details,
	}
}

func Unauthorized(message string) *AppError {
	return New(
		CodeUnauthorized,
		message,
		http.StatusUnauthorized,
	)
}

func Forbidden(message string) *AppError {
	return New(
		CodeForbidden,
		message,
		http.StatusForbidden,
	)
}

func NotFound(message string) *AppError {
	return New(
		CodeNotFound,
		message,
		http.StatusNotFound,
	)
}

func Conflict(message string) *AppError {
	return New(
		CodeConflict,
		message,
		http.StatusConflict,
	)
}

func TooManyRequests(message string) *AppError {
	return New(
		CodeTooManyRequests,
		message,
		http.StatusTooManyRequests,
	)
}

func Internal(message string, err error) *AppError {
	return &AppError{
		Code:       CodeInternalServer,
		Message:    message,
		HTTPStatus: http.StatusInternalServerError,
		Err:        err,
	}
}

func ServiceUnavailable(message string, err error) *AppError {
	return &AppError{
		Code:       CodeServiceUnavailable,
		Message:    message,
		HTTPStatus: http.StatusServiceUnavailable,
		Err:        err,
	}
}
