# /unslop - cut filler, keep substance

Applies to answers, docs, commit messages, PR descriptions and comments. Write in the user's language.

## Cut

- Openers and closers: greetings, praise, "here is...", "as mentioned above", offers of further help.
- Restating the question or the task.
- Hedging that carries no information ("it might possibly be worth considering").
- Explaining what the code already says.
- Adjectives and intensifiers with no measurable claim.

## Keep

- Names, numbers, versions, commands and paths, exactly as they are.
- Caveats that change what the reader will do.
- Decisions and the reason for them.
- Steps the reader must follow, in order.

## Method

1. Draft.
2. Delete everything that does not change what the reader knows or does.
3. Prefer one idea per sentence, active voice, concrete nouns.
4. Check that no fact, caveat or step was lost.

## For code

Delete dead code, comments that restate the code, wrappers that only forward calls, and unused parameters. Do not touch behavior.
