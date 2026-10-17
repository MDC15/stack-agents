# /verify - evidence before "done"

Runs automatically before you report a task as finished.

## Steps

1. **Acceptance list.** Re-read the original request. Write the acceptance points as a short checklist.
2. **Run checks.** Use the commands in "Project overrides" in `AGENTS.md`. If they are blank, discover them from `package.json`, `Makefile`, `pyproject.toml`, `Cargo.toml`, CI config or the README. Run, in order: format and lint, typecheck, targeted tests, then the wider suite, then build.
3. **Exercise the change.** Run the changed behavior directly: start the CLI, call the endpoint, execute the script, open the page. Passing tests alone do not cover UI, CLI or integration behavior.
4. **Review your own diff** (`git diff`, or compare against the originals if there is no git). Look for stray debug output, unrelated edits, secrets, leftover TODOs and commented-out code.
5. **Report** in two lists:
   - Verified: the command and its result.
   - Not verified: what, and why (no environment, no credentials, no tests).

## Rules

- Never say "done", "fixed" or "works" without a command result behind it.
- If a check fails, fix the cause. Do not weaken or skip the check.
- If a check was already failing before your change, show that with evidence and report it separately.
