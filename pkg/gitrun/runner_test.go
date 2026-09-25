package gitrun

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/behaviorengineering/majordomo-forge-clients/pkg/auth"
)

func TestAuthArgsDoNotContainRawToken(t *testing.T) {
	args := auth.GitConfigArgs(auth.Credential{Token: "sekrit", SCM: "github"})
	joined := strings.Join(args, " ")
	if strings.Contains(joined, "sekrit") {
		t.Fatal("token leaked into git config args")
	}
}

func TestIsRepoMissingDeadlineStillWorks(t *testing.T) {
	r := Runner{}
	if r.IsRepo("/tmp") {
		t.Skip("unexpected git repo in /tmp")
	}
}

func TestRunRequiresGit(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r := Runner{}
	_, err := r.RunTrim(ctx, "", "version")
	if err != nil {
		t.Fatalf("git version failed: %v", err)
	}
}
