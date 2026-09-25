package gitinspect

import (
	"context"
	"strings"

	"github.com/behaviorengineering/majordomo-forge-clients/pkg/auth"
	"github.com/behaviorengineering/majordomo-forge-clients/pkg/gitrun"
)

// Inspector performs read-only git inspection in a repository directory.
type Inspector struct {
	Run gitrun.Runner
}

// NewInspector returns an inspector with the given credentials.
func NewInspector(cred auth.Credential) *Inspector {
	return &Inspector{Run: gitrun.Runner{Cred: cred}}
}

// WorktreeStatus returns porcelain status lines (empty when clean).
func (in *Inspector) WorktreeStatus(ctx context.Context, dir string) ([]string, error) {
	out, err := in.Run.RunTrim(ctx, dir, "status", "--porcelain")
	if err != nil {
		return nil, err
	}
	if out == "" {
		return nil, nil
	}
	return strings.Split(out, "\n"), nil
}

// OriginURL returns the origin remote URL when present.
func (in *Inspector) OriginURL(ctx context.Context, dir string) (string, error) {
	return in.Run.RunTrim(ctx, dir, "remote", "get-url", "origin")
}
