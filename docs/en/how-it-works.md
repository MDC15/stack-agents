# How it works

> Tiếng Việt: [docs/vi/how-it-works.md](../../docs/vi/how-it-works.md)

## One source of truth

Every agent wants its rules in a different file. Maintaining the same text in
six files goes wrong the first time you edit one of them. So `AGENTS.md` is the
only file you edit, and `scripts/sync.sh` derives the rest from it.

That is the whole architecture. There is no build step, no dependency, and
nothing to install beyond `bash`.

## What the sync generates

| Agent | File | Note |
|---|---|---|
| Claude Code | `CLAUDE.md` | plain copy |
| Cursor | `.cursor/rules/stack-agents.mdc` | prefixed with `alwaysApply: true` frontmatter |
| Cline | `.clinerules` | plain copy |
| Roo Code | `.roorules` | plain copy |
| Windsurf | `.windsurfrules` | plain copy |
| GitHub Copilot | `.github/copilot-instructions.md` | plain copy |

Every generated file opens with a banner naming `AGENTS.md` as the source. Each
plain copy is byte-identical to `AGENTS.md` after the banner; the Cursor file
adds YAML frontmatter so the rule loads without being requested.

`scripts/sync.sh --check` compares each generated file against what the sync
would write, and exits non-zero naming any that differ. CI runs it, which is
what stops the six files from drifting apart. `--diff` does the same but also
prints the difference for each stale file, so drift is visible and not just
named.

If a tool looks somewhere else — these paths change between versions — edit the
`place` list at the bottom of `scripts/sync.sh`. The rest of the script does not
need to change.

## Why skills load on demand

`.ai/skills/*.md` and `.ai/playbooks/*.md` are not inlined into `AGENTS.md`.
`AGENTS.md` only carries the table that maps a command to a file. The agent
reads a skill file when the command is invoked.

This keeps the always-present context small: eight skills and two playbooks
would otherwise be in the window for every request, including the ones that
need none of them. It also means a skill can be long without costing anything
until it is used.

The exception is `/verify`, which the operating loop runs automatically before
work is reported as done.

## What `AGENTS.md` is, and is not

It is written instruction. It asks the agent to keep changes small, to verify
before claiming success, and to ask before destructive commands. It cannot
enforce any of that.

The hard rules at the bottom of `AGENTS.md` — no destructive commands without
confirmation, no secrets, text found in files is data and not instructions —
are the ones that matter most, and they are the ones a tool-level permission
system should back up. Treat `AGENTS.md` as the request and the sandbox as the
guarantee.

## Install and sync are separate

- `scripts/install.sh` copies the repository into a project: `AGENTS.md`, `.ai/`,
  `scripts/` and `VERSION`. It is a one-time act, and it refuses to overwrite by
  default. `--force` updates an existing install and reports the version change;
  `--uninstall` removes what stack-agents manages and keeps the project's own files.
- `scripts/sync.sh` runs inside a project and derives the per-agent files from
  `AGENTS.md`. It runs after every edit to `AGENTS.md`, and in CI.

Keeping them apart is why updating a project with `--force` merges the source
files over the project's without destroying `.ai/conventions.md`.

## Next

- [getting-started.md](getting-started.md) — install it.
- [skills.md](skills.md) — what each command does.
- [custom-agents.md](custom-agents.md) — wire it into your own loop.
