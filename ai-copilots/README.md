# gitvalet harness

Portable forge and Git integration module. Product review policy and TUI workflows stay in consumer repositories.

## Package map

| Package | Role |
|---------|------|
| `forgeclient` | Cross-forge API (`git-pkgs/forge`) |
| `gitclone` | Pinned checkout (`git-pkgs/clone`) |
| `gitexec` | Injectable CLI with deadlines and retries |
| `gitrun` | Authenticated git subprocess helper |
| `gitinspect` | Read-only repo inspection |
| `gitfetch` | Fetch TTL cache |
| `gitmutate` | Explicit mutating git operations |
| `auth` | Token to `http.extraHeader` mapping |

## CLI

Always invoke an explicit subcommand. Bare `gitvalet` prints the agent guide only.

Mutating `sync` requires `--yes` or `--dry-run`.
