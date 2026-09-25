package cli

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/behaviorengineering/gitvalet/pkg/auth"
	"github.com/behaviorengineering/gitvalet/pkg/gitclone"
	"github.com/behaviorengineering/gitvalet/pkg/gitinspect"
	"github.com/behaviorengineering/gitvalet/pkg/gitmutate"
	"github.com/behaviorengineering/gitvalet/pkg/gitrun"
)

// Exit codes for CI callers.
const (
	ExitOK           = 0
	ExitUsage        = 2
	ExitPrecondition = 3
	ExitRemote       = 4
	ExitAuth         = 5
	ExitCancel       = 6
	ExitInternal     = 1
)

var version = "dev"

// SetVersion injects the release identity from main.
func SetVersion(v string) {
	version = strings.TrimSpace(v)
	if version == "" {
		version = "dev"
	}
}

// Run dispatches argv to an explicit subcommand.
func Run(ctx context.Context, args []string, w io.Writer) int {
	if len(args) == 0 {
		printAgentGuide(w)
		return ExitOK
	}
	switch args[0] {
	case "help":
		printHelp(w)
		return ExitOK
	case "version":
		fmt.Fprintf(w, "gitvalet %s\n", version)
		return ExitOK
	case "inspect":
		return runInspect(ctx, args[1:], w)
	case "resolve":
		return runResolve(ctx, args[1:], w)
	case "checkout":
		return runCheckout(ctx, args[1:], w)
	case "sync":
		return runSync(ctx, args[1:], w)
	default:
		printHelp(w)
		return ExitUsage
	}
}

func printAgentGuide(w io.Writer) {
	fmt.Fprintln(w, "gitvalet: CI-friendly forge and git client")
	fmt.Fprintln(w, "Load ai-copilots/BOOTSTRAP.md for harness wiring.")
	fmt.Fprintln(w, "Commands: help, version, inspect, resolve, checkout, sync")
	fmt.Fprintln(w, "Mutating: sync requires --yes; use --dry-run to plan.")
}

func printHelp(w io.Writer) {
	fmt.Fprintln(w, "Usage: gitvalet <command> [flags]")
	fmt.Fprintln(w, "Commands:")
	fmt.Fprintln(w, "  help       human command catalog")
	fmt.Fprintln(w, "  version    release identity")
	fmt.Fprintln(w, "  inspect    read-only repository state")
	fmt.Fprintln(w, "  resolve    resolve remote branch to SHA")
	fmt.Fprintln(w, "  checkout   materialize pinned checkout")
	fmt.Fprintln(w, "  sync       fast-forward pull (mutating)")
}

func runInspect(ctx context.Context, args []string, w io.Writer) int {
	fs := flag.NewFlagSet("inspect", flag.ContinueOnError)
	dir := fs.String("dir", ".", "repository directory")
	jsonOut := fs.Bool("json", false, "JSON output")
	_ = fs.Parse(args)
	ctx, cancel := withDefaultDeadline(ctx)
	defer cancel()
	in := gitinspect.NewInspector(auth.Credential{})
	url, err := in.OriginURL(ctx, *dir)
	if err != nil {
		fmt.Fprintf(w, "inspect: %v\n", err)
		return ExitPrecondition
	}
	ref, ok := gitinspect.ParseRemoteURL(url)
	if *jsonOut {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"origin_url": url,
			"remote":     ref,
			"parsed":     ok,
		})
		return ExitOK
	}
	fmt.Fprintf(w, "origin=%s parsed=%v host=%s path=%s\n", url, ok, ref.Host, ref.Path)
	return ExitOK
}

func runResolve(ctx context.Context, args []string, w io.Writer) int {
	fs := flag.NewFlagSet("resolve", flag.ContinueOnError)
	remote := fs.String("remote", "", "HTTPS remote URL")
	branch := fs.String("branch", "main", "branch name")
	token := fs.String("token", "", "forge token (or FORGE_TOKEN env)")
	scm := fs.String("scm", "", "github|gitlab|bitbucket")
	jsonOut := fs.Bool("json", false, "JSON output")
	_ = fs.Parse(args)
	if strings.TrimSpace(*remote) == "" {
		return ExitUsage
	}
	cred := credentialFromFlags(*token, *scm, *remote)
	ctx, cancel := withDefaultDeadline(ctx)
	defer cancel()
	client := &gitclone.Client{Cred: cred}
	sha, err := client.ResolveBranchHead(ctx, *remote, *branch)
	if err != nil {
		fmt.Fprintf(w, "resolve: %v\n", err)
		return ExitRemote
	}
	if *jsonOut {
		_ = json.NewEncoder(w).Encode(map[string]string{"sha": sha, "branch": *branch})
		return ExitOK
	}
	fmt.Fprintln(w, sha)
	return ExitOK
}

func runCheckout(ctx context.Context, args []string, w io.Writer) int {
	fs := flag.NewFlagSet("checkout", flag.ContinueOnError)
	remote := fs.String("remote", "", "HTTPS remote URL")
	dst := fs.String("dst", "", "destination directory")
	sha := fs.String("sha", "", "pinned commit SHA")
	token := fs.String("token", "", "forge token (or FORGE_TOKEN env)")
	scm := fs.String("scm", "", "github|gitlab|bitbucket")
	dryRun := fs.Bool("dry-run", false, "plan only")
	_ = fs.Parse(args)
	if strings.TrimSpace(*remote) == "" || strings.TrimSpace(*dst) == "" || strings.TrimSpace(*sha) == "" {
		return ExitUsage
	}
	if *dryRun {
		fmt.Fprintf(w, "would checkout %s into %s\n", *sha, *dst)
		return ExitOK
	}
	cred := credentialFromFlags(*token, *scm, *remote)
	ctx, cancel := withDefaultDeadline(ctx)
	defer cancel()
	client := &gitclone.Client{Cred: cred}
	if err := client.EnsureCheckout(ctx, *remote, *dst, *sha); err != nil {
		fmt.Fprintf(w, "checkout: %v\n", err)
		return ExitRemote
	}
	fmt.Fprintf(w, "checked out %s\n", *sha)
	return ExitOK
}

func runSync(ctx context.Context, args []string, w io.Writer) int {
	fs := flag.NewFlagSet("sync", flag.ContinueOnError)
	dir := fs.String("dir", ".", "repository directory")
	branch := fs.String("branch", "main", "branch to fast-forward")
	yes := fs.Bool("yes", false, "confirm mutating sync")
	dryRun := fs.Bool("dry-run", false, "plan only")
	_ = fs.Parse(args)
	if !*yes && !*dryRun {
		fmt.Fprintln(w, "sync: mutating command requires --yes or --dry-run")
		return ExitUsage
	}
	if *dryRun {
		fmt.Fprintf(w, "would fast-forward %s on %s\n", *branch, *dir)
		return ExitOK
	}
	ctx, cancel := withDefaultDeadline(ctx)
	defer cancel()
	run := gitrun.Runner{}
	if err := gitmutate.FastForwardPull(ctx, run, *dir, *branch); err != nil {
		fmt.Fprintf(w, "sync: %v\n", err)
		if errors.Is(ctx.Err(), context.Canceled) {
			return ExitCancel
		}
		return ExitPrecondition
	}
	fmt.Fprintln(w, "synced")
	return ExitOK
}

func credentialFromFlags(token, scm, remote string) auth.Credential {
	if strings.TrimSpace(token) == "" {
		token = os.Getenv("FORGE_TOKEN")
	}
	if strings.TrimSpace(scm) == "" {
		scm = auth.InferSCM(remote)
	}
	return auth.Credential{Token: token, SCM: scm}
}

func withDefaultDeadline(ctx context.Context) (context.Context, context.CancelFunc) {
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Deadline(); ok {
		return ctx, func() {}
	}
	return context.WithTimeout(ctx, 2*time.Minute)
}
