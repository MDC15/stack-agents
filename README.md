# stack-agents

A set of Markdown rules and workflows that works with any AI coding agent:
Claude Code, Cursor, Codex, Cline, Roo Code, Windsurf, OpenCode, GitHub
Copilot, or an agent you build yourself with LangChain or CrewAI.

The idea: **one source of truth (`AGENTS.md`)**, rendered by a single static
binary (`stack`) into each tool's own rule file. Skills load on demand, so
they never sit in the context window unused.

**Languages:** English | [Tiếng Việt](README.vi.md)

## How it works

`AGENTS.md` holds the operating loop, the skill table and the hard rules. You
edit that file and nothing else. `stack sync` renders it into a `.agents/`
directory and links or copies it to every fixed path your agent reads —
nothing else to run, nothing to drift.

Skills and playbooks live under `.ai/` and load on demand: the agent reads
`.ai/skills/tdd.md` when `/tdd` is invoked, and not before.

More: [docs/en/how-it-works.md](docs/en/how-it-works.md).

## Structure

```
AGENTS.md                       # Single source of truth: loop, skill table, hard rules
VERSION                         # Current version; embedded in the binary
.ai/
├── principles.md               # 16 engineering principles (read for non-trivial work)
├── skills/                     # Invoked as /...
│   ├── architect.md            # Design before writing code
│   ├── arena.md                # Compare approaches before committing
│   ├── red-team.md             # Adversarial review (alias /swarm)
│   ├── tdd.md                  # Red - green - refactor
│   ├── verify.md               # Evidence before "done" (runs automatically)
│   ├── checkpoint.md           # Compare state around risky commands
│   ├── unslop.md               # Cut filler, keep substance
│   └── conventions.md          # Learn this repo's habits, save to .ai/conventions.md
└── playbooks/                  # Chosen by /auto <goal>
    ├── debug.md                # Includes a production-incident branch
    └── refactor.md
cmd/stack/                      # CLI: init, status, sync, check, uninstall, version
internal/
├── agents/                     # Registry: 8 agents and the fixed paths they read
├── detect/                     # Finds installed agents on this machine
├── generate/                   # Renders AGENTS.md into each file format
├── project/                    # Install, sync, check and uninstall one project
└── ui/                         # Agent table and the y/n prompt
embed.go                        # Embeds AGENTS.md, VERSION, .ai/ for shipping
.agents/                        # Generated; committed so CI can check drift
scripts/load_skill.py           # Loader for custom agents
docs/en/                        # Canonical documentation (English)
docs/vi/                        # Vietnamese translation of docs/en
CHANGELOG.md                    # What changed, per version
CONTRIBUTING.md                 # Workflow and pull-request checklist
Makefile                        # make check | sync | diff | test
```

## Install

`stack` is a single binary. Three ways to get it:

```bash
go install github.com/MDC15/stack-agents/cmd/stack@latest
```

or download an archive from the GitHub Releases page (linux, darwin,
windows × amd64, arm64), or build from a checkout:

```bash
git clone https://github.com/MDC15/stack-agents && cd stack-agents && go build -o stack ./cmd/stack
```

Then, in any project:

```bash
stack init --path /path/to/your-project
```

This detects the agents installed on the machine, shows the plan, asks once,
and writes `AGENTS.md`, `VERSION`, `.ai/`, `scripts/load_skill.py` plus the
generated `.agents/` files and the per-agent links. It refuses to touch an
existing install without `--force`, so a second run cannot destroy anything.
A `--force` update merges the new `AGENTS.md` over yours while keeping your
**Project overrides** section.

Real output, on a machine with OpenCode installed:

```
╭──────────────┬─────────────────┬─────────────────┬───────────────────────────────╮
│AGENT         │VERSION         │STATUS           │TARGET                         │
├──────────────┼─────────────────┼─────────────────┼───────────────────────────────┤
│Claude Code   │—               │✓ link           │CLAUDE.md                      │
│OpenCode      │opencode v2.0.20│✓ reads AGENTS.md│                               │
│Codex         │—               │✓ reads AGENTS.md│                               │
│Cursor        │—               │✓ link           │.cursor/rules/stack-agents.mdc │
│GitHub Copilot│—               │✗ skip           │.github/copilot-instructions.md│
│Cline         │—               │✓ link           │.clinerules                    │
│Roo Code      │—               │✗ skip           │.roorules                      │
│Windsurf      │—               │✗ skip           │.windsurfrules                 │
╰──────────────┴─────────────────┴─────────────────┴───────────────────────────────╯
AGENTS.md
.ai/playbooks/debug.md
... (one line per created file, then the links)
Next: edit AGENTS.md, then run `stack sync --path .`
```

| Option | Effect |
|---|---|
| `--agents <csv>` | Link exactly these agent IDs (claude, opencode, codex, cursor, copilot, cline, roo, windsurf) |
| `--all` | Link every agent |
| `--force` | Update an existing install, preserving Project overrides |
| `--dry-run` | Print the plan, change nothing |
| `--yes` | Skip the confirmation prompt |
| `--path` | Project directory (default `.`) |

> **Do not edit the shipped skill files** in an installed project — a
> `--force` update replaces them. Add your own `.ai/skills/<name>.md` instead;
> extra files are never touched by an update or an uninstall.

Full walkthrough: [docs/en/getting-started.md](docs/en/getting-started.md).

## Generated files

| Agent | Generated under `.agents/` | The agent reads | How it is linked |
|---|---|---|---|
| Claude Code | `claude.md` | `CLAUDE.md` | symlink |
| OpenCode | — | `AGENTS.md` | direct, no file |
| Codex | — | `AGENTS.md` | direct, no file |
| Cursor | `cursor.mdc` | `.cursor/rules/stack-agents.mdc` | symlink |
| GitHub Copilot | `copilot.md` | `.github/copilot-instructions.md` | copy with banner |
| Cline | `cline.md` | `.clinerules` | symlink |
| Roo Code | `roo.md` | `.roorules` | symlink |
| Windsurf | `windsurf.md` | `.windsurfrules` | symlink |

Presence is opt-in by default: only detected agents get links, plus the ones
you name with `--agents` or `--all`. After every edit to `AGENTS.md`, run
`stack sync` again. `stack check` fails naming any stale file — that is what
CI runs, so the tool-specific files cannot drift apart.

Adding a new agent means one row in
[`internal/agents/agents.go`](internal/agents/agents.go): see
[docs/en/extending.md](docs/en/extending.md).

## Commands

| Command | Effect |
|---|---|
| `stack init` | Detect agents, show the plan, install and link |
| `stack status` | Per-agent state: ok / drift / missing / skipped |
| `stack sync` | Re-render `.agents/*` and refresh the planned links |
| `stack check` | Fail naming every stale file; `--diff` shows the differences |
| `stack uninstall` | Remove what stack installed, keep your own files |
| `stack version` | Print the embedded version |

Exit codes: 0 success, 1 drift or failure, 2 usage error. Every command and
flag: [docs/en/cli.md](docs/en/cli.md).

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
`AGENTS.md`, and running `stack sync`. Playbooks work the same way under
`.ai/playbooks/`. Keep `AGENTS.md` short; longer rules belong in `.ai/`, where
they load only when needed.

More: [docs/en/extending.md](docs/en/extending.md).

## Verification

```bash
make check     # gofmt, go vet, tests, then `stack check`
```

The test suite covers the generated-file render and drift detection, install
and uninstall behaviour, agent detection, the CLI itself (subprocess tests
against a built binary), the `load_skill.py` CLI, every relative link in the
documentation, and that `docs/en` and `docs/vi` stay in step. CI runs the same
suite plus the generated-file drift check and Ruff over the Python loader.

## Contributing

Workflow, versioning and the pull-request checklist:
[CONTRIBUTING.md](CONTRIBUTING.md). What changed recently:
[CHANGELOG.md](CHANGELOG.md).

## Limitations

- Detection is best-effort: `stack` looks for each agent's binary on PATH,
  a marker directory in `$HOME`, or an installed VS Code extension. Agents
  installed another way are reported as not found; name them with
  `--agents` explicitly.
- Symlinks are rewritten to copies (with the same banner) when the OS refuses
  to create one, so an install still works end to end.
- Symlinked rule files render poorly on a Windows git checkout: they either
  become plain text files or need developer mode enabled. The copy fallback
  in the previous bullet is the same story.
- These are written instructions, not an enforcement mechanism. How closely an
  agent follows them depends on the model and the tool. Dangerous operations
  still need permissions and a sandbox on the tool's side.

## License

[MIT](LICENSE).