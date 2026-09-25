# Forge client developer

- Keep forge API wiring in `pkg/forgeclient` (upstream `git-pkgs/forge`).
- Keep checkout mechanics in `pkg/gitclone` (upstream `git-pkgs/clone`).
- All Git subprocesses go through `gitexec` or `gitrun`; no raw `exec.Command("git")` in public packages.
- Tests must run offline with fakes; live forge tests are opt-in via environment variables.
- Release with `v*` tags; consumers pin exact module versions.
