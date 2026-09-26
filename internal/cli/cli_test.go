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
	s := buf.String()
	if !strings.Contains(s, "gitvalet") {
		t.Fatal("expected guide")
	}
	if !strings.Contains(s, "AGENTS.md") {
		t.Fatal("expected AGENTS.md in guide")
	}
}

func TestHelpAliases(t *testing.T) {
	for _, arg := range []string{"help", "-h", "--help"} {
		var buf bytes.Buffer
		if Run(context.Background(), []string{arg}, &buf) != ExitOK {
			t.Fatalf("%s failed", arg)
		}
		if !strings.Contains(buf.String(), "Usage:") {
			t.Fatalf("%s missing usage", arg)
		}
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
