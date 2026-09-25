package main

import (
	"context"
	"os"
	"runtime/debug"

	"github.com/behaviorengineering/gitvalet/internal/cli"
)

var buildVersion string

func main() {
	cli.SetVersion(version())
	os.Exit(cli.Run(context.Background(), os.Args[1:], os.Stdout))
}

func version() string {
	if buildVersion != "" && buildVersion != "dev" {
		return buildVersion
	}
	if info, ok := debug.ReadBuildInfo(); ok {
		if info.Main.Version != "" && info.Main.Version != "(devel)" {
			return info.Main.Version
		}
	}
	return "dev"
}
