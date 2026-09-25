package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestBareInvokePrintsGuide(t *testing.T) {
	var buf bytes.Buffer
	code := Run(context.Background(), nil, &buf)
	if code != ExitOK {
		t.Fatalf("code %d", code)
	}
	if !strings.Contains(buf.String(), "majordomo-forge") {
		t.Fatal("expected guide")
	}
}

func TestUnknownCommand(t *testing.T) {
	var buf bytes.Buffer
	code := Run(context.Background(), []string{"nope"}, &buf)
	if code != ExitUsage {
		t.Fatalf("code %d", code)
	}
}

func TestVersion(t *testing.T) {
	SetVersion("v0.1.0-test")
	var buf bytes.Buffer
	code := Run(context.Background(), []string{"version"}, &buf)
	if code != ExitOK || !strings.Contains(buf.String(), "v0.1.0-test") {
		t.Fatalf("unexpected: %q code=%d", buf.String(), code)
	}
}

func TestSyncDryRun(t *testing.T) {
	var buf bytes.Buffer
	code := Run(context.Background(), []string{"sync", "--dry-run", "--branch", "main"}, &buf)
	if code != ExitOK {
		t.Fatalf("code %d", code)
	}
}
