package domain

import (
	"errors"
	"fmt"
	"testing"
)

func TestAppError_ErrorWithUnderlying(t *testing.T) {
	underlying := fmt.Errorf("db connection failed")
	appErr := &AppError{Code: "DB_ERR", Message: "database error", Err: underlying}

	want := "database error: db connection failed"
	if got := appErr.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestAppError_ErrorWithoutUnderlying(t *testing.T) {
	appErr := &AppError{Code: "BAD_REQ", Message: "bad request", Err: nil}

	want := "bad request"
	if got := appErr.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestAppError_Unwrap(t *testing.T) {
	underlying := fmt.Errorf("wrapped")
	appErr := &AppError{Code: "X", Message: "msg", Err: underlying}

	if got := appErr.Unwrap(); got != underlying {
		t.Errorf("Unwrap() = %v, want %v", got, underlying)
	}
}

func TestNewAppError(t *testing.T) {
	underlying := ErrNotFound
	appErr := NewAppError("GAME_NOT_FOUND", "game not found", underlying)

	if appErr.Code != "GAME_NOT_FOUND" {
		t.Errorf("Code = %q, want %q", appErr.Code, "GAME_NOT_FOUND")
	}
	if appErr.Message != "game not found" {
		t.Errorf("Message = %q, want %q", appErr.Message, "game not found")
	}
	if appErr.Err != underlying {
		t.Errorf("Err = %v, want %v", appErr.Err, underlying)
	}
}

func TestAppError_ErrorsIs(t *testing.T) {
	appErr := NewAppError("NOT_FOUND", "resource missing", ErrNotFound)

	if !errors.Is(appErr, ErrNotFound) {
		t.Error("errors.Is(appErr, ErrNotFound) = false, want true")
	}
}
