# Playbook: refactor

Change structure, not behavior. Never mix the two in one step.

## Steps

1. **Goal.** State what is being improved and which behavior must not change.
2. **Safety net.** Check test coverage of the code you will touch. Where it is missing, add characterization tests that pin current behavior, including quirks. Commit them before refactoring.
3. **Plan.** Break the work into small steps, each of which leaves the build and tests green. For legacy code, prefer introducing a seam and moving callers over one step at a time to a big-bang rewrite.
4. **Execute.** For each step: change, run tests, checkpoint (see `.ai/skills/checkpoint.md` for bulk changes), suggest a commit.
5. **Stop condition.** If tests go red and the cause is not clear within one attempt, revert that step and report.
6. **Bugs you notice.** Write them down and report them separately. Do not fix them inside the refactor unless the user agrees.
7. **Verify.** Follow `.ai/skills/verify.md`.
8. **Report.** What moved, what stayed identical, what you found along the way.
