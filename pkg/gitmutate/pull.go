package gitmutate

import (
	"context"
	"fmt"
	"strings"

	"github.com/behaviorengineering/gitvalet/pkg/gitrun"
)

// FastForwardPull runs git pull --ff-only for one branch when the worktree is clean.
func FastForwardPull(ctx context.Context, run gitrun.Runner, repoDir, branch string) error {
	branch = strings.TrimSpace(branch)
	if branch == "" {
		return fmt.Errorf("gitmutate: branch required")
	}
	status, err := run.RunTrim(ctx, repoDir, "status", "--porcelain")
	if err != nil {
		return err
	}
	if strings.TrimSpace(status) != "" {
		return fmt.Errorf("gitmutate: refusing pull on dirty worktree")
	}
	_, err = run.Run(ctx, repoDir, "pull", "--ff-only", "origin", branch)
	return err
}
