// Package forgeclient adapts git-pkgs/forge for majordomo consumers.
package forgeclient

import (
	"strings"

	forge "github.com/git-pkgs/forge"
)

// RepositoryRef is a validated forge repository identity.
type RepositoryRef struct {
	Domain string
	Owner  string
	Name   string
}

// ParseRepositoryURL parses a forge HTTPS clone URL into a RepositoryRef.
func ParseRepositoryURL(raw string) (RepositoryRef, error) {
	domain, owner, name, err := forge.ParseRepoURL(strings.TrimSpace(raw))
	if err != nil {
		return RepositoryRef{}, err
	}
	return RepositoryRef{Domain: domain, Owner: owner, Name: name}, nil
}

// NewForgeClient constructs the upstream cross-forge API client.
func NewForgeClient(opts ...forge.Option) *forge.Client {
	return forge.NewClient(opts...)
}
