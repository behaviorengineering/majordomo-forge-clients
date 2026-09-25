# gitvalet

A git and forge attendant for Go applications: shared forge API and Git checkout client library for majordomo, gitboard, and other hosts.

- `pkg/forgeclient`: cross-forge API adapter around `git-pkgs/forge`
- `pkg/gitclone`: pinned checkout adapter around `git-pkgs/clone`
- `pkg/gitexec`: deadline-gated CLI execution with retries
- `pkg/auth`: HTTPS token header policy for Git smart HTTP
- `cmd/gitvalet`: CI-friendly CLI (`inspect`, `resolve`, `checkout`, `sync`)

See `AGENTS.md` for agent operating rules.

## Releases (for agents)

Default branch merges that change releasable code get patch tags via CI. Consumers pin `github.com/behaviorengineering/gitvalet@vX.Y.Z`.
