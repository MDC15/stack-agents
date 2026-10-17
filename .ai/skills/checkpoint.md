# /checkpoint - compare state around risky commands

Use around commands that change many files or are hard to undo: migrations, codegen, bulk rewrites, dependency upgrades, formatters over the whole repo. Skip it for read-only commands.

## Before

1. Record state: `git status --short` and `git rev-parse HEAD`. Outside git, list the relevant directory with modification times.
2. If the tree has uncommitted work that the command could overwrite, tell the user and propose a WIP commit or `git stash`. Do it only after they agree.

## After

1. Run `git status --short` and `git diff --stat` again.
2. Compare with the "before" state. Confirm that only the files you expected changed.

## If something unexpected changed

Stop. Report which files and how they differ. Do not run further destructive commands to clean up. Offer options (revert specific files, restore from the checkpoint) and wait.
