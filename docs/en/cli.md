# CLI reference

> Tiếng Việt: [docs/vi/cli.md](../../docs/vi/cli.md)

`stack` has six commands. Exit codes are consistent everywhere: 0 success,
1 drift or failure, 2 usage error (unknown command, bad flag).

```
usage: stack <command> [args]
commands: --version, check, init, status, sync, uninstall, version
```

## stack init

Escalates a directory into a stack-agents project: detect, show the plan as a
table, confirm (unless `--yes`), copy the ship set (`AGENTS.md`, `VERSION`,
`.ai/**`, `scripts/load_skill.py`), render `.agents/*`, and link or copy the
fixed path of every planned agent.

| Flag | Command(s) | Effect |
|---|---|---|
| `--path <dir>` | all | Project directory, default `.` |
| `--yes` | init, uninstall | Skip the confirmation prompt |
| `--dry-run` | init, sync, uninstall | Print the plan, change nothing |
| `--force` | init | Update an existing install, merge preserving Project overrides |
| `--agents <csv>` | init, sync | Link these agent IDs (adds to what was detected) |
| `--all` | init, sync | Link every agent |

Detection is presence-as-opt-in: agents found on the machine get links
automatically, everything else is skipped. An unknown agent ID is an error
(exit 1).

## stack status

Reads the project and reports every agent: `ok`, `drift`, `missing link`, or
`skipped (agent not found)`. Exits 1 when the directory is not a project.

## stack sync

Re-renders every `.agents/*` file from the current `AGENTS.md`, then creates
or refreshes the fixed-path entry for each planned agent. A stack-owned
symlink is re-pointed; a stack-owned copy is rewritten; an entry that never
came from stack is left alone and reported. An agent's link stays put when
the agent is unplanned — presence is opt-in. An edit to `AGENTS.md` makes the
generated files stale until the next sync; CI runs [the check](#stack-check)
to catch that.

## stack check

Compares every `.agents/*` file and fixed-path entry against what the sync
would write right now. Each stale entry prints one line:

```
out of date: .agents/claude.md
```

`--diff` appends a unified-style diff for each stale `.agents/*` file:
current lines prefixed `-`, expected lines prefixed `+`. Exit 1 when anything
stale or missing; also exits 1 when a detected agent's link is gone, or when
the directory is not a project.

## stack uninstall

Removes the `.agents/` directory, the fixed-path entries that are symlinks
into it or banner-carrying copies, `VERSION`, `scripts/load_skill.py`, and
every `.ai/**` file stack shipped. It keeps `AGENTS.md`, `.ai/conventions.md`,
and every skill you added yourself. Without `--yes` it confirms first;
`--dry-run` prints the list instead of removing.

## stack version

Prints the version embedded in the binary, `stack 1.0.0` (which is what
`--version` does; the alias is kept for usability).

## Migration from ≤ 0.2

Version 0.2 shipped bash scripts (`install.sh`, `sync.sh`, `test.sh`); they
are gone in 1.0. To move an existing project over:

1. Still running on ≤ 0.2 tools? Run the old uninstall first:
   `bash /path/to/stack-agents/scripts/install.sh --uninstall <project>`.
   This is optional: a fresh `stack init` converts the layout in place.
2. Run `stack init --path <project>`. A fresh install or `--force` replaces
   what remains, and root rule files that still carry the old
   `scripts/sync.sh` banner are converted automatically — the same first
   `stack sync` recognizes both banner generations as stack-owned.
3. Run `stack check` and add it to CI so the generated files cannot drift.

## Next

- [getting-started.md](getting-started.md) — install and verify.
- [how-it-works.md](how-it-works.md) — why `.agents/` and symlinks.
- [extending.md](extending.md) — add an agent, a skill, a playbook.