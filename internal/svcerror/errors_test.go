package svcerror

import (
	"errors"
	"net/http"
	"testing"
)

func TestValidationError_Error(t *testing.T) {
	err := &ValidationError{Message: "title is required"}
	if got := err.Error(); got != "title is required" {
		t.Errorf("Error() = %q, want %q", got, "title is required")
	}
}

func TestValidationError_Unwrap(t *testing.T) {
	err := &ValidationError{Message: "invalid"}
	if !errors.Is(err, ErrValidation) {
		t.Error("ValidationError should unwrap to ErrValidation")
	}
}

func TestHTTPStatusAndBody(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode int
		wantErr  string
		wantMsg  string
	}{
		{
			name:     "nil error",
			err:      nil,
			wantCode: http.StatusOK,
			wantErr:  "",
			wantMsg:  "",
		},
		{
			name:     "ErrTodoNotFound",
			err:      ErrTodoNotFound,
			wantCode: http.StatusNotFound,
			wantErr:  "not_found",
			wantMsg:  "todo not found",
		},
		{
			name:     "ValidationError with message",
			err:      &ValidationError{Message: "title is required"},
			wantCode: http.StatusBadRequest,
			wantErr:  "validation_error",
			wantMsg:  "title is required",
		},
		{
			name:     "ValidationError empty message",
			err:      &ValidationError{Message: ""},
			wantCode: http.StatusBadRequest,
			wantErr:  "validation_error",
			wantMsg:  "validation failed",
		},
		{
			name:     "unknown error",
			err:      errors.New("database connection failed"),
			wantCode: http.StatusInternalServerError,
			wantErr:  "internal_error",
			wantMsg:  "An error occurred",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body := HTTPStatusAndBody(tt.err)
			if code != tt.wantCode {
				t.Errorf("status = %d, want %d", code, tt.wantCode)
			}
			if tt.err == nil {
				if body != nil {
					t.Errorf("body = %v, want nil", body)
				}
				return
			}
			if body["error"] != tt.wantErr {
				t.Errorf("body[error] = %q, want %q", body["error"], tt.wantErr)
			}
			if body["message"] != tt.wantMsg {
				t.Errorf("body[message] = %q, want %q", body["message"], tt.wantMsg)
			}
		})
	}
}
