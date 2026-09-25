// Package auth maps forge credentials to safe Git HTTPS configuration.
package auth

import (
	"encoding/base64"
	"strings"
)

// Credential is a forge-scoped HTTPS token for Git smart HTTP.
type Credential struct {
	Token string
	SCM   string // github, gitlab, bitbucket, or empty for GitHub-style default
}

// GitConfigArgs returns git -c http.extraHeader=... args. Empty token yields nil.
func GitConfigArgs(cred Credential) []string {
	token := strings.TrimSpace(cred.Token)
	if token == "" {
		return nil
	}
	scm := strings.ToLower(strings.TrimSpace(cred.SCM))
	var header string
	switch scm {
	case "gitlab":
		basic := base64.StdEncoding.EncodeToString([]byte("oauth2:" + token))
		header = "Authorization: Basic " + basic
	case "bitbucket":
		basic := base64.StdEncoding.EncodeToString([]byte("x-token-auth:" + token))
		header = "Authorization: Basic " + basic
	default:
		basic := base64.StdEncoding.EncodeToString([]byte("x-access-token:" + token))
		header = "Authorization: Basic " + basic
	}
	return []string{"-c", "http.extraHeader=" + header}
}

// InferSCM guesses the forge from an HTTPS remote URL host.
func InferSCM(remoteURL string) string {
	u := strings.ToLower(remoteURL)
	switch {
	case strings.Contains(u, "gitlab"):
		return "gitlab"
	case strings.Contains(u, "bitbucket"):
		return "bitbucket"
	default:
		return "github"
	}
}

// RedactToken returns a safe placeholder when token is non-empty.
func RedactToken(token string) string {
	if strings.TrimSpace(token) == "" {
		return ""
	}
	return "<redacted>"
}
