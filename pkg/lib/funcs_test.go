package lib

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
	"time"
)

func TestMaskedText(t *testing.T) {
	testCases := []struct {
		input    string
		n        int
		expected string
	}{
		{"1234567890", 3, "123...890"},
		{"hello world", 5, "hel...rld"},
		{"test", 10, "t...t"},
	}

	for _, tc := range testCases {
		result := MaskedText(tc.input, tc.n)
		if result != tc.expected {
			t.Errorf("MaskedText(%s, %d) = %s; want %s", tc.input, tc.n, result, tc.expected)
		}
	}
}

func TestIsDeadlineExceeded(t *testing.T) {
	// тест: ошибка DeadlineExceeded в grpc.Status
	deadlineErr := status.Error(codes.DeadlineExceeded, "context deadline exceeded")
	if !IsDeadlineExceeded(deadlineErr) {
		t.Errorf("Expected true, got false for DeadlineExceeded error")
	}

	// тест: обычная ошибка
	normalErr := errors.New("random error")
	if IsDeadlineExceeded(normalErr) {
		t.Errorf("Expected false, got true for non-DeadlineExceeded error")
	}

	// тест: DeadlineExceeded без grpc.Status
	deadlineExceededErr := context.DeadlineExceeded
	if !IsDeadlineExceeded(deadlineExceededErr) {
		t.Errorf("Expected true, got false for clear DeadlineExceeded error")
	}
	// v2
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	time.Sleep(1 * time.Second)
	if !IsDeadlineExceeded(ctx.Err()) {
		t.Errorf("Expected true, got false for error: %v", ctx.Err())
	}

}
