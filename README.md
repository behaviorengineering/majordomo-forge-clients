# Majordomo Forge Clients

Shared forge API and Git checkout client library for majordomo and gitboard.

- `pkg/forgeclient`: cross-forge API adapter around `git-pkgs/forge`
- `pkg/gitclone`: pinned checkout adapter around `git-pkgs/clone`
- `pkg/gitexec`: deadline-gated CLI execution with retries
- `pkg/auth`: HTTPS token header policy for Git smart HTTP
- `cmd/majordomo-forge`: CI-friendly CLI (`inspect`, `resolve`, `checkout`, `sync`)

See `AGENTS.md` for agent operating rules.
