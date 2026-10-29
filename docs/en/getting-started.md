# Getting started

> Tiếng Việt: [docs/vi/getting-started.md](../../docs/vi/getting-started.md)

## Prerequisites

- `stack` — one static binary. Install with
  `go install github.com/MDC15/stack-agents/cmd/stack@latest`, from a release
  archive, or by building a checkout: `go build -o stack ./cmd/stack`.
- `python3`, only for custom agents (`scripts/load_skill.py`) and one repo
  test; the CLI itself needs nothing else.

## Initialize a project

```bash
stack init --path /path/to/your-project
```

The output is a table of the eight agents found or skipped on this machine.

```
╭──────────────┬─────────────────┬─────────────────╮
│ AGENT        │ VERSION         │ STATUS          │
├──────────────┼─────────────────┼─────────────────┤
│ Claude Code  │ —               │ ✓ link          │
│ OpenCode     │ opencode v2.0.2 │ ✓ reads AGENTS. │
│ Codex        │ —               │ ✗ skip          │
│ Cursor       │ —               │ ✗ skip          │
│ GitHub Copilot│ —              │ ✓ link          │
│ Cline        │ —               │ ✗ skip          │
│ Roo Code     │ —               │ ✗ skip          │
│ Windsurf     │ —               │ ✗ skip          │
╰──────────────┴─────────────────┴─────────────────╯
```

Then one confirmation (`--yes` skips it) and the install: the ship set is
copied, `.agents/` is rendered, and every planned agent's fixed path appears.

```
Next: edit AGENTS.md, then run `stack sync --path /path/to/your-project`
```

On the second run, without `--force`, it refuses with
"already installed" — nothing is overwritten by accident.

| Option | Effect |
|---|---|
| `--agents <csv>` | Link these agents even when detection missed them |
| `--all` | Link every agent |
| `--force` | Update an existing install; merges the new `AGENTS.md` and keeps your **Project overrides** |
| `--dry-run` | Print the plan, change nothing |
| `--yes` | Skip the prompt |

> **Do not edit the shipped skill files** in `.ai/skills/` and
> `.ai/playbooks/` — a `--force` update replaces them. Add your own
> `.ai/skills/<name>.md` instead; extra files survive updates and uninstalls
> untouched.

## Fill in Project overrides

`AGENTS.md` ends with a **Project overrides** section. These take priority
over everything above them except the hard rules. Put the real commands for
this project there:

```
- Build: `npm run build`
- Test: `npm test`
- Lint/typecheck: `npm run lint`
```

`/verify` reads this section to decide what to run. If it is left blank, the
agent has to guess from `package.json`, `Makefile`, `pyproject.toml`,
`Cargo.toml` or CI config.

## Regenerate after every edit

Edit `AGENTS.md`, never the generated files. Then:

```bash
stack sync --path /path/to/your-project
```

To check without writing — in CI, for example:

```bash
stack check --path /path/to/your-project          # exits 1 and names any stale file
stack check --diff --path /path/to/your-project   # same, and prints the difference
```

## Verify the install

```bash
stack status --path /path/to/your-project
stack check --path /path/to/your-project
```

`status` shows one row per agent; `check` should print "no drift". Confirm the
agent picked up the config by asking for something small, for example
`/verify`.

Never worked? See the current known caveat in
[cli.md](cli.md#migration-from-≤-02) — before `stack` existed, projects were
installed by bash scripts; a first `stack sync` recognizes that era by the
legacy banner and converts the files.

## What your agent now sees

Each tool reads its own file, generated from `.agents/`: `CLAUDE.md` for
Claude Code, `.cursor/rules/stack-agents.mdc` for Cursor, `.clinerules` for
Cline, `.roorules` for Roo Code, `.windsurfrules` for Windsurf,
`.github/copilot-instructions.md` for Copilot, and `AGENTS.md` itself for
Codex and OpenCode. The full table, with notes, is in
[how-it-works.md](how-it-works.md#generated-files).

## Next

- [how-it-works.md](how-it-works.md) — what the sync actually generates.
- [skills.md](skills.md) — the commands, one by one.
- [extending.md](extending.md) — add your own skill or playbook.
- [../../README.md](../../README.md) — the short version.