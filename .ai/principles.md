# Engineering principles

## Context
1. **Narrow context.** Load the smallest set of files that answers the question. Search before reading; read ranges before whole files.
2. **Navigate once, query many.** Map a directory once, then reason from the map instead of re-reading.
3. **Ground claims.** Do not state that a function, flag or API exists unless you saw it in the repo or docs this session.

## Code
4. **Match the codebase.** Follow existing naming, structure, error handling and test style over personal preference.
5. **Smallest sufficient change.** Prefer editing to adding, adding to rewriting. No speculative abstractions.
6. **Comments explain why, not what.** Write a comment only for non-obvious intent, constraints, workarounds or public API contracts. Delete comments that restate code. Keep license headers and required doc comments.
7. **Make invalid states unrepresentable.** Use types, enums and constructors before runtime checks.
8. **Fail loudly at boundaries.** Validate input where it enters the system. Do not swallow errors; handle them or propagate them with context.
9. **No dead code.** Remove unused code you introduce. Leave no commented-out code or stray debug output.
10. **Dependencies are liabilities.** Add one only when writing it yourself is clearly worse. State why.

## Process
11. **Evidence over confidence.** "Should work" is not a result. Run it, or say you could not.
12. **Reproduce before fixing.** A bug without a reproduction is a guess.
13. **Reversible steps.** Prefer changes that are easy to undo. Checkpoint before risky work.
14. **Surface uncertainty.** If you assumed something, say what and why. One line is enough.

## Communication
15. **Lead with the result.** Answer first, detail after.
16. **No filler.** No praise, no restating the question, no announcing what you are about to do.
