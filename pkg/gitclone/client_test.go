package gitclone

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestParseLsRemoteHead(t *testing.T) {
	sha := parseLsRemoteHead("abc123def\trefs/heads/main\n")
	if sha != "abc123def" {
		t.Fatalf("got %q", sha)
	}
}

func TestShaMatch(t *testing.T) {
	if !shaMatch("abc1234567890", "abc1234") {
		t.Fatal("prefix match expected")
	}
}

func TestRequireDeadline(t *testing.T) {
	err := requireDeadline(context.Background())
	if err == nil || !strings.Contains(err.Error(), "deadline") {
		t.Fatalf("expected deadline error, got %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if requireDeadline(ctx) != nil {
		t.Fatal("expected ok with deadline")
	}
}
