# Playbook: debug

One variable at a time. No shotgun changes.

## Steps

1. **Reproduce.** Get the exact error, input and environment. Turn it into a failing test or a script that fails reliably. If you cannot reproduce it, say so and list what you need.
2. **Isolate.** Shrink the failing case. Use `git log` or `git bisect` when it used to work. Write two or three hypotheses, ranked by likelihood, and test the cheapest first with evidence (logs, a debugger, a minimal probe).
3. **Root cause.** State in one sentence why it happens, not just where it shows up.
4. **Fix.** Make the smallest change that removes the cause. Add or keep the regression test from step 1.
5. **Clean up.** Revert failed experiments. Remove temporary logging.
6. **Verify.** Follow `.ai/skills/verify.md`.
7. **Report.** Cause, fix, regression test, residual risk.

## If it is a production incident

- Start read-only: logs, metrics, recent deploys and config changes. Change nothing in production without explicit confirmation.
- Propose mitigation first (rollback, feature flag, scale up) and the cheapest way to confirm it worked. A proper fix comes after.
- Record a short timeline: what was observed, when, and what was changed.
