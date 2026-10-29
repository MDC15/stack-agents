# Contributing

## The one rule

`AGENTS.md` is the single source of truth. Edit it, never the files it
generates: `CLAUDE.md`, `.clinerules`, `.roorules`, `.windsurfrules`,
`.cursor/rules/stack-agents.mdc`, `.github/copilot-instructions.md`. After every
edit:

```bash
make sync
```

Commit the regenerated files together with the `AGENTS.md` change, so the
drift check in CI stays green.

## Adding a skill or playbook

1. Create `.ai/skills/<name>.md` (playbooks go under `.ai/playbooks/`).
   Names may contain only lowercase letters, digits and hyphens.
2. Register it in the skill table in `AGENTS.md`, so every agent can find it.
3. Run `make sync`.

Full walkthrough: [docs/en/extending.md](docs/en/extending.md).

## Versioning

Bump `VERSION` and add a `CHANGELOG.md` entry when a change reaches someone
installing or extending stack-agents. Fixes to prose alone need neither.

## Documentation

English is canonical (`docs/en/`). Keep `docs/vi/` in step: same filenames,
same headings. A bilingual parity test fails when either side is missing a
file or drifts in heading count.

## Checklist before a pull request

1. `make check` passes — generated-file drift check plus the test suite. CI
   runs the same command.
2. `AGENTS.md` and the generated files are in sync.
3. New skills are registered in the `AGENTS.md` table.
4. Both language versions of a documentation change are present.
5. `VERSION` and `CHANGELOG.md` are updated when the change is user-visible.
