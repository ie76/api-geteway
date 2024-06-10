package errors

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	code := ErrInvalidInput
	expectedMessage := "Invalid input"
	expectedRetryable := false

	err := New(code)

	if err.Code != code {
		t.Errorf("Expected error code %s, got %s", code, err.Code)
	}

	if err.Message != expectedMessage {
		t.Errorf("Expected error message %s, got %s", expectedMessage, err.Message)
	}

	if err.Retryable != expectedRetryable {
		t.Errorf("Expected error retryable %t, got %t", expectedRetryable, err.Retryable)
	}

}

func TestRetry(t *testing.T) {
	mockFunction := func() error {
		return New(ErrDatabaseError)
	}

	err := Retry(3, time.Second, func() error {
		return New(ErrUnauthorized)
	})

	if err == nil {
		t.Errorf("Expected non-nil error for non-retryable error, got nil")
	}

	err = Retry(3, time.Second, mockFunction)

	if err == nil {
		t.Errorf("Expected non-nil error for retryable error, got nil")
	}
}
