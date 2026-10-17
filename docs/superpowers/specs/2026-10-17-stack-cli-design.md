# stack CLI — design

Date: 2026-10-17
Status: approved in conversation; builder check resolved (see "Resolved after review")

## Intent and why

stack-agents ships `AGENTS.md` as the single source of truth and two bash
scripts (`install.sh`, `sync.sh`) that generate per-agent rule files. Two
problems the user hit:

1. Install dumps every per-agent config file into the project root, even for
   agents the project does not use (`.windsurfrules` for someone who never
   uses Windsurf). Six generated files scatter across fixed paths
   (`CLAUDE.md`, `.clinerules`, `.roorules`, `.windsurfrules`,
   `.github/copilot-instructions.md`, `.cursor/rules/stack-agents.mdc`).
2. The agent list is hard-coded in bash; there is no way to see which agents
   the machine actually has, or what a sync produced.

The refactor: one CLI that detects the AI coding agents installed on the
machine, generates the configs they need grouped in a single `.agents/`
folder, and links only what the detected agents read. Project-scoped,
never global.

## Goals

- Group all generated per-agent configs into one folder: `.agents/`.
- Auto-detect installed AI coding agents; only detected agents get files at
  the fixed paths they read. No stray files for unused agents.
- One thin file at each fixed path (symlink into `.agents/`), so rule content
  exists in exactly one copy and cannot drift.
- Simple, beautiful CLI (Go + Bubble Tea + lipgloss) showing the needed
  info: detected agents, versions, link state, drift.
- Project-scoped: nothing written outside the target project.
- Repo cleanup: `install.sh` / `sync.sh` / `test.sh` go away, replaced by the
  CLI; this repo itself dogfoods the new layout.
- `AGENTS.md` stays the single source of truth; `.ai/` skills and playbooks
  are unchanged; `scripts/load_skill.py` stays for custom agents.

## Non-goals (v1)

- Per-agent rule customization. All agents get the same rules generated from
  `AGENTS.md`; Cursor additionally gets its frontmatter.
- Global install of rule files under the user's home — project only.
- Editing rule content inside the TUI. Content lives in `AGENTS.md`, edited
  in the user's editor.
- Supporting every editor/agent in existence. The eight below are v1;
  adding one later is one table row.
- Windows-first support (a copy fallback covers it; see below).

## Anti-goals

- Creating files for agents the user does not have installed.
- Duplicating rule content across many root files that silently drift.
- Requiring a runtime (Python, Node) to run the CLI itself. Single static
  binary.

## The decisions (made in conversation)

### 1. One folder, symlinks at fixed paths

Each agent reads its config at a hard-coded path and most have no import
mechanism (only Claude Code supports `@file` imports, in `CLAUDE.md`).
Stubs cannot work for Cursor/Copilot/Cline/Roo/Windsurf; symlinks work for
every agent and keep exactly one copy of the content. Layout
(user-approved):

```
project/
├── AGENTS.md               # single source of truth, edited by the user
├── .ai/                    # skills + playbooks, unchanged
├── VERSION
├── scripts/load_skill.py   # custom-agent loader, unchanged
└── .agents/                # ALL generated content lives here
    ├── claude.md           CLAUDE.md                          → .agents/claude.md
    ├── cursor.mdc          .cursor/rules/stack-agents.mdc     → .agents/cursor.mdc
    ├── copilot.md          .github/copilot-instructions.md    → .agents/copilot.md
    ├── cline.md            .clinerules                        → .agents/cline.md
    ├── roo.md              .roorules                          → .agents/roo.md
    └── windsurf.md         .windsurfrules                     → .agents/windsurf.md
```

### 2. Only detected agents get linked

| Agent | Detected by | Generated link |
|---|---|---|
| Claude Code | `claude` in PATH or `~/.claude/` exists | `CLAUDE.md` → `.agents/claude.md` |
| OpenCode | `opencode` in PATH or `~/.config/opencode/` exists | none — reads `AGENTS.md` directly |
| Codex | `codex` in PATH or `~/.codex/` exists | none — reads `AGENTS.md` directly |
| Cursor | `cursor` in PATH or `~/.cursor/` exists | `.cursor/rules/stack-agents.mdc` → `.agents/cursor.mdc` |
| GitHub Copilot | VS Code extension `github.copilot` | `.github/copilot-instructions.md` → `.agents/copilot.md` |
| Cline | VS Code extension `saoudrizwan.claude-dev` | `.clinerules` → `.agents/cline.md` |
| Roo Code | VS Code extension `rooveterinaryinc.roo-cline` | `.roorules` → `.agents/roo.md` |
| Windsurf | `windsurf` in PATH or `~/.codeium/` exists | `.windsurfrules` → `.agents/windsurf.md` |

Detection is best-effort: PATH lookup plus per-user config directories;
Cline/Roo/Copilot ride on `code --list-extensions` (skipped cleanly when the
`code` binary is absent). Version, when cheaply available (`<tool>
--version`), is shown in the table.

Links are opt-in by presence: a fixed-path file that already exists as one
of our symlinks (or banner-carrying copy) is kept and re-synced even if
detection misses it. `--agents <list>` and `--all` create links regardless
of detection — also the remedy for detection false negatives. This is how
this repo's own six committed links get created: one maintainer run of
`stack init --all`.

### 3. CLI: Go + Bubble Tea + lipgloss, single static binary named `stack`

`go:embed` embeds what an install ships (`AGENTS.md`, `.ai/`, `VERSION`,
`scripts/load_skill.py`), so `stack init` works from a bare binary — no
clone of this repo needed.

Commands:

| Command | Effect |
|---|---|
| `stack init [path]` | Detect agents, show table, confirm, copy shipped files, generate `.agents/`, link; print next steps. Refuses an existing install without `--force` (merge semantics as today: project's `.ai/conventions.md` and custom skills survive). |
| `stack status` | One table: agent, version, detected?, link state (ok / missing / drift / skipped). |
| `stack sync` | Regenerate `.agents/*` from `AGENTS.md`; create or refresh links for detected agents. |
| `stack check [--diff]` | Exit 1 when a generated file or link drifts from `AGENTS.md` (for CI). |
| `stack uninstall` | Remove `.agents/`, the links it created (verified as symlinks into `.agents/` or banner-carrying copies), and shipped files. Keep `AGENTS.md`, `.ai/conventions.md`, custom skills. |
| `stack version` | Print version. |

Shared flags: `--yes` (skip confirmation, for scripts), `--dry-run`
(print what would happen; init, uninstall, sync).

UI: lipgloss tables, `✓`/`✗` status columns, plain-language messages.
CLI messages are **English** (decided with the user after review: the repo's
canonical language; docs stay bilingual). `init` and `status` render the
same table component. Mockup shape shown to and approved by the user
(mockup was drafted in Vietnamese; only the layout was approved):

```
$ stack init
  Detecting AI coding agents...

  Agent          Version   Status
  ✓ Claude Code  2.4.1     create  .agents/claude.md   → CLAUDE.md
  ✓ Cursor       1.9.8     create  .agents/cursor.mdc  → .cursor/rules/…
  ✓ Copilot      —         create  .agents/copilot.md  → .github/…
  ✗ Windsurf     not found skip
  ✗ Roo Code     not found skip

  Continue? [Y/n]
```

### 4. Fallback: copy instead of symlink

Where symlink creation fails (Windows checkout without developer mode, FAT
mounts), the CLI writes a copy of the `.agents/` file at the fixed path
carrying the generated banner; `stack check` and `stack sync` compare copies
against their source the same way. No hard dependency on symlink support.

### 5. Generated file format

Unchanged semantics from `sync.sh`: `.agents/<agent>.md` = banner comment +
`AGENTS.md` content; `.agents/cursor.mdc` = Cursor frontmatter
(`description`, `alwaysApply: true`) + banner + content. The banner names
`stack sync` as the generator.

### 6. Repo dogfood

- The repo's own root files become symlinks into `.agents/`, committed as
  symlinks in git.
- `.github/workflows/check.yml`: gofmt/vet check, `go test ./...`, Ruff over
  `load_skill.py`, then `go run ./cmd/stack check`.
- `scripts/test.sh` is removed. Its repo-level checks (relative-link
  integrity across published docs, docs/en ↔ docs/vi parity) move into Go
  tests. CLI behavior tests (generate, detect, init/uninstall, drift) are Go
  tests too, with fake HOME/PATH.

### 7. Version and distribution

- Bump `VERSION` to 1.0.0 (breaking: scripts gone, layout changes).
- Distribution: GitHub Releases with prebuilt binaries (linux/darwin/windows,
  amd64/arm64) plus `go install`. Release workflow details are the
  builder's call.

### 8. Docs

- `README.md` + `README.vi.md` rewritten around the CLI.
- `docs/en/*` and `docs/vi/*` updated; a `cli.md` page documents the
  commands. Bilingual parity stays enforced by tests.
- Migration note for existing installs: run the old
  `scripts/install.sh --uninstall`, then `stack init`.

## Constraints

- Go toolchain for development; end users need only the binary.
- Bilingual docs parity stays enforced (same filenames, same heading counts).
- Relative-link integrity over published docs stays enforced.
- `scripts/load_skill.py` keeps its path and behavior (custom agents
  reference it).
- `stack init` must preserve the "Project overrides" semantics: the user
  fills in build/test/lint commands at the end of their `AGENTS.md`; sync
  regenerates `.agents/*` from the whole file, overrides included.

## Resolved after review

Builder check raised six questions. One went to the user (CLI language =
English). The rest, settled here:

1. **Shipped template** — `stack init` embeds and copies this repo's own
   `AGENTS.md` verbatim (same as `install.sh` does today), kept current
   with the new toolchain so no shipped instruction references a file that
   no longer exists.
2. **`--force` preserves Project overrides** — on update, everything before
   the `## Project overrides` heading is replaced with the upstream
   version; the target's section from that heading on is kept (today
   `install.sh --force` silently wipes it — that is a bug, not a
   semantic to keep).
3. **Missing links fail `check` only when detected** — drift, broken
   symlinks, wrong link targets, and banner-carrying copies out of sync
   always fail; an absent link fails only when its agent is detected on the
   machine running the check. CI runners detect nothing, so the CI verdict
   is deterministic (content + existing-link integrity).
4. **`init` adopts a pre-existing `AGENTS.md`** — if the target has an
   `AGENTS.md` stack did not ship (no `VERSION` marker from a prior
   install), it becomes the source of truth untouched; `init` only adds
   `.agents/`, links, and the shipped support files around it.
5. **Old-layout migration is automatic** — root files carrying the
   `Generated from AGENTS.md` banner are recognized as ours and replaced
   with the new layout on `init`/`sync`; no manual uninstall round-trip.
6. **Copilot file is a copy, not a symlink** — github.com's server-side
   Copilot reads the git blob, where a symlink is stored as a path string.
   `.github/copilot-instructions.md` is therefore always a banner-carrying
   copy; all other fixed paths stay symlinks.

## What's left to the builder

- Go module layout (suggested: `cmd/stack`, `internal/detect`,
  `internal/generate`, `internal/install`, `internal/ui`), naming, error
  handling.
- Exact Bubble Tea/lipgloss styling.
- Detection details: timeouts for `--version` calls, behavior when neither
  PATH nor config dir exists, `code` binary absent.
- Confirmation flow details; `--yes` and `--dry-run` plumbing.
- Release workflow shape.
- Test layout: faking HOME/PATH for detection tests; temp-dir integration
  tests for init/sync/check/uninstall.
