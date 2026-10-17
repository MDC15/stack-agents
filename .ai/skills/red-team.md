# /red-team (alias /swarm) - adversarial review

Attack a design or a diff before it ships. Do not defend it.

## Lenses

- **Security:** input validation, authn/authz, injection, secrets, data exposure, dependency risk.
- **Operations:** failure and retry behavior, timeouts, load, observability, rollback, migrations, config drift.
- **QA:** untested paths, edge cases, concurrency, bad input, backward compatibility, ambiguous requirements.

## Steps

1. Identify the artifact under review (design doc, files, or `git diff`) and read it fully.
2. Run each lens in turn. If your environment supports subagents, run each lens in a separate subagent that sees only the artifact and the requirements, so the lenses do not anchor on each other. Otherwise run them sequentially and say that independence is limited.
3. Per lens, report at most five findings. Each needs: severity (high, medium, low), evidence (file and line, or a concrete scenario), and a suggested fix. Drop speculative findings you cannot tie to evidence.
4. Close with a verdict: **ship**, **fix first**, or **redesign**.
5. Do not apply fixes automatically. Ask which ones to apply.

## Output

A single table of findings, then the verdict.
