package svcerror

import (
	"errors"
	"net/http"
)

// Sentinel errors for the service layer.
var (
	ErrTodoNotFound = errors.New("todo not found")
	ErrValidation   = errors.New("validation error")
)

// ValidationError wraps a validation message.
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func (e *ValidationError) Unwrap() error {
	return ErrValidation
}

// HTTPStatusAndBody maps service errors to HTTP status and JSON body.
// Returns (statusCode, body map for {"error": "...", "message": "..."}).
func HTTPStatusAndBody(err error) (int, map[string]string) {
	if err == nil {
		return http.StatusOK, nil
	}
	if errors.Is(err, ErrTodoNotFound) {
		return http.StatusNotFound, map[string]string{
			"error":   "not_found",
			"message": "todo not found",
		}
	}
	if errors.Is(err, ErrValidation) {
		msg := err.Error()
		if msg == "" || msg == ErrValidation.Error() {
			msg = "validation failed"
		}
		return http.StatusBadRequest, map[string]string{
			"error":   "validation_error",
			"message": msg,
		}
	}
	return http.StatusInternalServerError, map[string]string{
		"error":   "internal_error",
		"message": "An error occurred",
	}
}
