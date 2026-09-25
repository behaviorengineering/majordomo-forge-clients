# Bootstrap

Resolve the module directory:

```bash
go list -m -f '{{.Dir}}' github.com/behaviorengineering/gitvalet
```

Link consumer IDE skills to `ai-copilots/skills/` in that directory (symlink or junction). Do not copy skill bodies into host `.cursor/skills/`.
