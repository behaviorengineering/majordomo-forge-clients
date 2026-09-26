# BOOTSTRAP — gitvalet ai-copilots

**Audience:** Any AI agent (Cursor, GitHub Copilot, Claude Code, Codex) in a
workspace that depends on or checks out this module.

**Goal:** Wire host IDE discovery to canonical content under `ai-copilots/`.
Optionally refresh content. **MUST NOT** copy skill bodies unless symlinks or
junctions fail and the user approves copy fallback.

**Module path:** `github.com/behaviorengineering/gitvalet`

---

## When to run

| Mode | Phases |
|------|--------|
| **Wire only** | 0 → 2 → 3 → 4 |
| **Refresh content + wire** | 0 → 1 → 2 → 3 → 4 |

---

## Phase 0 — Resolve module root

From a Go module that requires `github.com/behaviorengineering/gitvalet` (or this checkout):

```bash
MOD="$(go list -m -f '{{.Dir}}' github.com/behaviorengineering/gitvalet)"
test -d "$MOD/ai-copilots" || { echo "missing ai-copilots under $MOD"; exit 1; }
echo "Module Dir: $MOD"
```

If `go list` is unavailable, use a known nested checkout path only when it
clearly contains `ai-copilots/`. MUST NOT invent a path.

Re-run wire after module version bumps (cached Dir can change).

---

## Phase 1 — Refresh content (optional)

Edit only files under `$MOD/ai-copilots/`. Load author-ai-copilots + agent-smith
when authoring skills.

---

## Phase 2 — Ask IDE and OS if unknown

1. IDE: Cursor, GitHub Copilot, Claude Code, Codex
2. OS: macOS/Linux symlink vs Windows junction/copy
3. Workspace: library alone vs nested under a parent monorepo vs dependency-only

---

## Phase 3 — Wire discovery

Canonical sources:

| Artifact | Path under `$MOD` |
|----------|-------------------|
| Skill tree | `ai-copilots/skills/gitvalet-operator/` |
| Skill tree | `ai-copilots/skills/gitvalet-developer/` |

Discovery paths:

| IDE | Skills |
|-----|--------|
| Cursor | `.cursor/skills/**/SKILL.md` |
| GitHub Copilot | `.github/skills/**/SKILL.md` |
| Claude Code | `.claude/skills/**/SKILL.md` |
| Codex | `.codex/skills/**/SKILL.md` |

**Cursor example (macOS/Linux)** from the **host workspace root**:

```bash
MOD="$(go list -m -f '{{.Dir}}' github.com/behaviorengineering/gitvalet)"
mkdir -p .cursor/skills
ln -snf "$MOD/ai-copilots/skills/gitvalet-operator" .cursor/skills/gitvalet-operator
ln -snf "$MOD/ai-copilots/skills/gitvalet-developer" .cursor/skills/gitvalet-developer
```

When the workspace root is this repository, relative links are fine:

```bash
ln -snf "$(pwd)/ai-copilots/skills/gitvalet-operator" .cursor/skills/gitvalet-operator
ln -snf "$(pwd)/ai-copilots/skills/gitvalet-developer" .cursor/skills/gitvalet-developer
```

**GitHub Copilot example:**

```bash
MOD="$(go list -m -f '{{.Dir}}' github.com/behaviorengineering/gitvalet)"
mkdir -p .github/skills
ln -snf "$MOD/ai-copilots/skills/gitvalet-operator" .github/skills/gitvalet-operator
ln -snf "$MOD/ai-copilots/skills/gitvalet-developer" .github/skills/gitvalet-developer
```

---

## Phase 4 — Verify

```bash
test -f .cursor/skills/gitvalet-operator/SKILL.md || test -f .github/skills/gitvalet-operator/SKILL.md
```

At least one IDE skill path MUST resolve after wire.
