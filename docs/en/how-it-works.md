# How it works

> Tiếng Việt: [docs/vi/how-it-works.md](../../docs/vi/how-it-works.md)

## One source of truth

Every agent wants its rules in a different file. Maintaining the same text in
six files goes wrong the first time you edit one of them. So `AGENTS.md` is
the only file you edit, and the `stack` binary derives the rest from it.

That is the whole architecture. One static binary, no runtime dependency,
nothing to leave on a machine but the files it generates.

## `.agents/` and the fixed paths

`stack sync` renders `AGENTS.md` once per output format into `.agents/`:

| Agent | Generated | The agent reads | How |
|---|---|---|---|
| Claude Code | `.agents/claude.md` | `CLAUDE.md` | relative symlink |
| OpenCode | — | `AGENTS.md` | reads it directly |
| Codex | — | `AGENTS.md` | reads it directly |
| Cursor | `.agents/cursor.mdc` | `.cursor/rules/stack-agents.mdc` | relative symlink |
| GitHub Copilot | `.agents/copilot.md` | `.github/copilot-instructions.md` | banner copy |
| Cline | `.agents/cline.md` | `.clinerules` | relative symlink |
| Roo Code | `.agents/roo.md` | `.roorules` | relative symlink |
| Windsurf | `.agents/windsurf.md` | `.windsurfrules` | relative symlink |

Symlinks point inward (`CLAUDE.md` → `.agents/claude.md`) so one rewrite
reaches every agent at once. Copilot does not follow them, so it gets a
banner-carrying copy rewritten on the spot. When the OS refuses a symlink,
`stack` falls back to the same banner copy, so an install always completes.

Every generated file opens with a banner naming `AGENTS.md` as the source:
`<!-- Generated from AGENTS.md by stack sync. Edit AGENTS.md, not this file. -->`
That banner is also the ownership marker: sync, check and uninstall only
replace a file that carries it, a symlink into `.agents/`, or a file the
embedded ship set enumerates. A file you own — a hand-written `.clinerules`,
your `.ai/conventions.md`, a custom skill — is never touched.

## Presence is opt-in

`stack sync` only touches the fixed paths of agents that are either detected
on this machine or named with `--agents` / `--all`. Everything else stays
away: an unplanned `.roorules` a project wrote by hand is left alone, even if
its name is in the registry. `stack status` shows the split; `--agents` adds
specifically; `--all` takes everything.

## Drift check

`stack check` compares each generated file and fixed path against what the
sync would write right now, and exits 1 naming every stale entry:

```
out of date: .agents/claude.md
```

An edit to `AGENTS.md` makes the generated files stale until the next sync;
CI runs the check, which is what stops the tool-specific files from drifting
apart. `--diff` prints the difference for each stale `.agents/*` file, so
drift is visible and not just named.

If a tool looks somewhere else — these paths change between versions — add a
row to the registry in
[`internal/agents/agents.go`](../../internal/agents/agents.go) and rerun
`make check`. Nothing else needs to change.

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

## Install, sync and uninstall are separate

- `stack init` makes a project a stack-agents project: it copies the ship
  set, renders `.agents/`, creates links. One-time; refuses to overwrite
  without `--force`, and `--force` merges while preserving **Project
  overrides**.
- `stack sync` runs inside a project and derives the per-agent files from
  `AGENTS.md`. It runs after every edit to `AGENTS.md`, and in CI.
- `stack uninstall` removes everything stack owns and keeps the project's own
  files, so leaving the system is a one-command act too.

Keeping them apart is why updating a project with `--force` merges the source
files over the project's without destroying `.ai/conventions.md`.

## Next

- [getting-started.md](getting-started.md) — install it.
- [cli.md](cli.md) — every command and flag.
- [skills.md](skills.md) — what each command does.
- [custom-agents.md](custom-agents.md) — wire it into your own loop.