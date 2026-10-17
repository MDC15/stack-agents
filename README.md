# pstack-universal

A set of Markdown rules and workflows that works with any AI coding agent:
Claude Code, Cursor, Codex, Cline, Roo Code, Windsurf, GitHub Copilot, or an
agent you build yourself with LangChain or CrewAI.

The idea: **one source of truth (`AGENTS.md`)**, with a script that generates
each tool's own rule file from it. Skills are read only when they are needed,
so they never sit in the context window unused.

**Languages:** English | [Tiếng Việt](README.vi.md)

## How it works

`AGENTS.md` holds the operating loop, the skill table and the hard rules. You
edit that file and nothing else. `scripts/sync.sh` writes its contents into
every per-tool rule file your agent looks for, so they cannot drift apart.

Skills and playbooks live under `.ai/` and load on demand: the agent reads
`.ai/skills/tdd.md` when `/tdd` is invoked, and not before.

More: [docs/en/how-it-works.md](docs/en/how-it-works.md).

## Structure

```
AGENTS.md                  # Single source of truth: loop, skill table, hard rules
VERSION                    # Current version; copied into every install
.ai/
├── principles.md          # 16 engineering principles (read for non-trivial work)
├── skills/                # Invoked as /...
│   ├── architect.md       # Design before writing code
│   ├── arena.md           # Compare approaches before committing
│   ├── red-team.md        # Adversarial review (alias /swarm)
│   ├── tdd.md             # Red - green - refactor
│   ├── verify.md          # Evidence before "done" (runs automatically)
│   ├── checkpoint.md      # Compare state around risky commands
│   ├── unslop.md          # Cut filler, keep substance
│   └── conventions.md     # Learn this repo's habits, save to .ai/conventions.md
└── playbooks/             # Chosen by /auto <goal>
    ├── debug.md           # Includes a production-incident branch
    └── refactor.md
scripts/
├── install.sh             # Install, update or uninstall with one command
├── sync.sh                # Generate each agent's rule file from AGENTS.md
├── test.sh                # Self-check: install.sh, sync.sh, load_skill.py, links
└── load_skill.py          # Loader for custom agents
docs/
├── en/                    # Canonical documentation (English)
└── vi/                    # Vietnamese translation of docs/en
CHANGELOG.md               # What changed, per version
CONTRIBUTING.md            # Workflow and pull-request checklist
Makefile                   # make check | sync | diff | test | install
```

## Install

One command, no manual copying:

```bash
bash /path/to/pstack-universal/scripts/install.sh /path/to/your-project
```

It copies `AGENTS.md`, `.ai/`, `scripts/` and `VERSION` into the project and
runs `sync.sh` there. Then fill in the build, test and lint commands under
**Project overrides** at the end of `AGENTS.md`.

| Option | Effect |
|---|---|
| `--force` | Merge the current `AGENTS.md`, `.ai/`, `scripts/` and `VERSION` over the target, then re-run the sync. This is the update path; it prints the old and new version |
| `--dry-run` | Print what would happen, change nothing |
| `--no-sync` | Copy the files but do not generate the agent rule files |
| `--uninstall` | Remove what pstack installed (generated files, scripts, shipped skills), keeping the project's own files. Honours `--dry-run` |
| `--version` | Print the version of the checkout and exit |

Without `--force`, the script refuses to touch an existing install, so it cannot
destroy work by accident. A merge replaces files the source contains and leaves
the rest alone; `.ai/conventions.md` and any skill you added yourself are not in
the source, so they survive an update.

> **Do not edit the shipped skill files** in an installed project — a
> `--force` update replaces them. Add your own `.ai/skills/<name>.md` instead;
> extra files are never touched by an update or an uninstall.

To remove pstack from a project:

```bash
bash /path/to/pstack-universal/scripts/install.sh --uninstall /path/to/your-project
```

<details>
<summary>Manual install, if the script does not fit your setup</summary>

1. Copy `AGENTS.md`, `.ai/` and `scripts/` into the project root.
2. Fill in the build, test and lint commands under **Project overrides** at
   the end of `AGENTS.md`.
3. Run `bash scripts/sync.sh`.

</details>

Full walkthrough: [docs/en/getting-started.md](docs/en/getting-started.md).

## Generated files

| Agent | File it reads |
|---|---|
| Claude Code | `CLAUDE.md` |
| Cursor | `.cursor/rules/pstack.mdc` (always applied) |
| Cline | `.clinerules` |
| Roo Code | `.roorules` |
| Windsurf | `.windsurfrules` |
| GitHub Copilot | `.github/copilot-instructions.md` |
| Codex and any agent that reads `AGENTS.md` | `AGENTS.md` directly |

Tool versions change these paths. If your agent looks somewhere else, edit the
`place` list in `scripts/sync.sh`.

After every edit to `AGENTS.md`, run `sync.sh` again. To see what a sync would
change without writing, run `bash scripts/sync.sh --diff` (or `make diff`).
In CI, run `bash scripts/sync.sh --check` to fail when a generated file has
drifted.

## Commands

| Command | Effect |
|---|---|
| `/architect` | Design (data, flow, interfaces), then stop and wait for approval |
| `/arena` | Two or three genuinely different approaches, a comparison table, a recommendation |
| `/red-team` | Attack the work through three lenses: security, operations, QA |
| `/tdd` | Failing test first, minimal code, then refactor |
| `/verify` | Run lint, tests, build and exercise the change. Runs automatically before "done" |
| `/checkpoint` | Record and compare state around a risky command |
| `/unslop` | Tighten prose, keeping every number and caveat |
| `/conventions` | Learn this repo's and this user's habits, save them after you approve |
| `/auto <goal>` | Pick the matching playbook (`debug`, `refactor`) and run it |

Each one in full: [docs/en/skills.md](docs/en/skills.md).

## Custom agents

For your own loop (LangChain, CrewAI, anything else), `scripts/load_skill.py`
exposes the rules as a system prompt and the skills as a tool:

```python
from scripts.load_skill import system_prompt, load_skill, list_modules

prompt = system_prompt()   # use as the agent's system prompt
# expose load_skill(name, kind="skills") as a tool so the agent loads skills on demand
```

`load_skill` accepts only `[a-z0-9-]` names and the kinds `skills` and
`playbooks`, so an agent cannot use it to read files outside `.ai/`.

More: [docs/en/custom-agents.md](docs/en/custom-agents.md).

## Extending

Add a skill by creating `.ai/skills/<name>.md`, adding one row to the table in
`AGENTS.md`, and running `sync.sh`. Playbooks work the same way under
`.ai/playbooks/`. Keep `AGENTS.md` short; longer rules belong in `.ai/`, where
they load only when needed.

More: [docs/en/extending.md](docs/en/extending.md).

## Verification

```bash
make check     # generated-file drift check, then the test suite
```

`make check` runs `scripts/sync.sh --check` and `scripts/test.sh`. The test
suite covers `install.sh` behaviour, `sync.sh` drift detection, the
`load_skill.py` CLI, that every relative link in the documentation resolves,
and that `docs/en` and `docs/vi` stay in step. CI runs the same command, plus
ShellCheck over the shell scripts and Ruff over the Python loader.

## Contributing

Workflow, versioning and the pull-request checklist:
[CONTRIBUTING.md](CONTRIBUTING.md). What changed recently:
[CHANGELOG.md](CHANGELOG.md).

## Limitations

- These are written instructions, not an enforcement mechanism. How closely an
  agent follows them depends on the model and the tool. Dangerous operations
  still need permissions and a sandbox on the tool's side.
- The per-agent file paths above are not verified against every current
  version of every tool. Try one simple request (for example `/verify`) to
  confirm your agent picked up the config.

## License

[MIT](LICENSE).
