# /conventions - learn this repo's habits

Capture conventions into `.ai/conventions.md` so future sessions follow them without being told again.

You cannot see past chats unless they are in your context. Use only evidence you can actually inspect.

## Sources

- Config: linter, formatter, tsconfig or equivalent, `.editorconfig`, CI workflows.
- Code: naming, file layout, error handling, test style in several representative files.
- History: `git log` messages for commit style, if available.
- This conversation: corrections and preferences the user has stated.

## Steps

1. Collect candidate rules from the sources above. Keep only those seen in at least two places or stated by the user.
2. Write each rule as one imperative line with its evidence in brackets, for example `Use named exports [seen in src/*, eslint rule]`.
3. Keep the file under 40 lines. Group by topic: naming, structure, errors, tests, commits.
4. Show the proposed file to the user. Save it only after they approve.

## Rules

- Never include secrets, personal data or anything from private messages.
- Mark a rule `(inferred)` when the user has not confirmed it.
- When asked to update, merge into the existing file instead of rewriting it.
