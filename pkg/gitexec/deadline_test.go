package gitexec

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRunRequiresDeadline(t *testing.T) {
	r := New()
	_, err := r.Run(context.Background(), "echo", "hi")
	if !errors.Is(err, ErrMissingDeadline) {
		t.Fatalf("expected ErrMissingDeadline, got %v", err)
	}
}

func TestRunWithDeadlineUsesFakeAttempt(t *testing.T) {
	r := New()
	r.attempt = func(ctx context.Context, _ time.Duration, name string, args ...string) ([]byte, error) {
		return []byte("ok"), nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	out, err := r.Run(ctx, "echo", "hi")
	if err != nil || string(out) != "ok" {
		t.Fatalf("unexpected: out=%q err=%v", out, err)
	}
}
