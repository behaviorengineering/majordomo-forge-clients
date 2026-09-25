# Forge client operator

1. Run `majordomo-forge help` for the command catalog.
2. Use `inspect` and `resolve` before mutating checkouts.
3. Pass tokens via `FORGE_TOKEN` or `--token`; never embed tokens in URLs.
4. Pin SHAs with `resolve` once, then `checkout --sha` for deterministic CI.
5. Use `sync --dry-run` then `sync --yes` for fast-forward pulls.

Exit codes: 0 success, 2 usage, 3 precondition, 4 remote, 5 auth, 6 cancel.
