# Extending

> Tiếng Việt: [docs/vi/extending.md](../../docs/vi/extending.md)

## Add a skill

1. Create `.ai/skills/<name>.md`. The name may contain only lowercase letters,
   digits and hyphens — `load_skill` rejects anything else.
2. Add one row to the skill table in `AGENTS.md`, mapping the command to the
   file and saying when to use it.
3. Run `stack sync` so the rendered files and links pick up the new table.

That is all. The skill is now available as `/<name>` in every tool, and through
`load_skill("<name>")` in a custom agent.

Write the skill file itself the way the existing ones are written: a short title,
when to use it, the steps, and what it produces. Keep it about the process, not
about a specific project.

## Add a playbook

Same shape, under `.ai/playbooks/`. Add the name to the `/auto` line in
`AGENTS.md` so the model knows it exists, and run `stack sync`.

## Add an agent

One row in `internal/agents/agents.go` — the ID, the display name, the kind
(symbolic link, banner copy, or direct like OpenCode reading `AGENTS.md`
itself), the file it is generated to under `.agents/`, and the fixed path the
tool reads. Detection signals live beside the registry in
`internal/detect/detect.go`; add the probe the agent needs — binary on PATH,
a marker directory under `$HOME`, or an installed VS Code extension. Commit
with `stack sync --all` so the generated files and links for the new agent
appear in the same commit; `make check` verifies both.

## Keep `AGENTS.md` short

`AGENTS.md` is in the context window for every request. Anything longer than a
couple of lines belongs in `.ai/` and loads only when needed. The skill table is
the right size for `AGENTS.md`; the skill's steps are not.

The rule of thumb: if the text only matters when a command is invoked, it belongs
in the file that command loads.

## Verify your change

```bash
make check
```

That runs gofmt, `go vet`, `go test ./...` (which includes the documentation
link and parity checks, the drift detection, install and uninstall behaviour,
the CLI subprocess tests, and the `load_skill.py` CLI) and finally
`stack check`. If you added a skill, confirm it shows up:

```bash
python3 -I scripts/load_skill.py list
```

## Keep the translations in step

Documentation exists twice — `docs/en/` and `docs/vi/`, with identical filenames.
If you change one, change its counterpart in the same commit. The bilingual
parity test fails when a file exists on one side only, or when the two versions
of a file carry different numbers of headings — both signs of a translation left
behind. Only stale prose slips through, so treat the two files as one change.

`AGENTS.md` and `.ai/` are deliberately English-only: they are loaded into the
model's context every session, and a second copy would double that cost.

## Contribute

The full checklist lives in [../../CONTRIBUTING.md](../../CONTRIBUTING.md).
For extending specifically:

1. `make check` passes.
2. `AGENTS.md` and the generated files are in sync — CI fails otherwise.
3. New skills are registered in the `AGENTS.md` table.
4. Both language versions of any documentation change are present.

## Next

- [skills.md](skills.md) — the existing skills, as a model to follow.
- [how-it-works.md](how-it-works.md) — why the split is drawn this way.
- [cli.md](cli.md) — the commands the checks ride on.
- [../../README.md](../../README.md) — the short version.