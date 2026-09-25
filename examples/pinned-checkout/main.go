// Package main is a compile-tested example of pinned checkout usage.
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/behaviorengineering/gitvalet/pkg/auth"
	"github.com/behaviorengineering/gitvalet/pkg/gitclone"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Fprintln(os.Stderr, "usage: example <remote-url> <branch> <dst>")
		os.Exit(2)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	client := &gitclone.Client{Cred: auth.Credential{Token: os.Getenv("FORGE_TOKEN"), SCM: auth.InferSCM(os.Args[1])}}
	sha, err := client.ResolveBranchHead(ctx, os.Args[1], os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := client.EnsureCheckout(ctx, os.Args[1], os.Args[3], sha); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(sha)
}
