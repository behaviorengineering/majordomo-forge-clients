---
name: gitvalet-operator
description: >-
  Operate the gitvalet CLI for inspect, resolve, checkout, and sync against
  GitHub, GitLab, or Bitbucket. Use when scripting forge checkouts, CI git
  steps, or debugging gitvalet exit codes.
---

# gitvalet operator

## Constraints

**CONSTRAINT:** MUST run `gitvalet help` before the first mutating command in a session.
- Enforcement: Operator checklist step 1
- Violation: STOP, run help, confirm command catalog

**CONSTRAINT:** MUST use read-only commands (`inspect`, `resolve`) before `checkout` or `sync`.
- Enforcement: Review command order in scripts and CI
- Violation: STOP, run inspect/resolve, then mutate

**CONSTRAINT:** MUST pass forge tokens via `FORGE_TOKEN` or `--token`; MUST NOT embed tokens in remote URLs.
- Enforcement: Grep scripts for `https://.*@` or token query params
- Violation: STOP, move token to env or flag

**CONSTRAINT:** MUST pin SHAs with `resolve` once, then `checkout --sha` for deterministic CI.
- Enforcement: CI logs show resolve output SHA reused in checkout
- Violation: STOP, add resolve step before checkout

**CONSTRAINT:** MUST use `sync --dry-run` before `sync --yes` for fast-forward pulls.
- Enforcement: Dry-run appears in logs or plan output before mutating sync
- Violation: STOP, run dry-run, then `--yes`

CORRECT:
```bash
export FORGE_TOKEN="$CI_JOB_TOKEN"
gitvalet resolve --remote "$HTTPS_URL" --branch main
gitvalet checkout --remote "$HTTPS_URL" --dst ./work --sha "$SHA"
gitvalet sync --dry-run --dir ./work --branch main
gitvalet sync --yes --dir ./work --branch main
```

PROHIBITED:
```bash
gitvalet checkout --remote "https://token@github.com/org/repo.git" ...
gitvalet sync --dir . --branch main   # missing --yes and --dry-run
```

## Exit codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 2 | Usage |
| 3 | Precondition (local repo/state) |
| 4 | Remote / network |
| 5 | Auth |
| 6 | Cancelled |

## Pre-completion verification

- [ ] **Catalog:** `gitvalet help` lists inspect, resolve, checkout, sync
      Method: Run help
      Pass: All commands listed
      Fail: STOP, upgrade gitvalet binary
- [ ] **Mutating gate:** sync without `--yes` fails usage unless `--dry-run`
      Method: `gitvalet sync` with no flags
      Pass: Exit 2
      Fail: STOP, fix script flags
