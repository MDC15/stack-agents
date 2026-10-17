# Changelog

Versions follow [semver](https://semver.org/). An entry lists what changed
for someone installing or extending stack-agents; prose-only fixes need
no entry.

## 0.2.0 - 2026-10-17

Renamed the project from `pstack-universal` to `stack-agents`.

Changed:

- The generated Cursor rule file is now `.cursor/rules/stack-agents.mdc`
  (was `.cursor/rules/pstack.mdc`). An upgrade via `install.sh --force` writes
  the new file but does not delete the old one; if a project still has
  `.cursor/rules/pstack.mdc`, remove it by hand so Cursor applies the rules
  once instead of twice.

## 0.1.0 - 2026-10-17

First versioned release; until now the repository carried no version number.

Added:

- `VERSION` file. `install.sh` copies it into the project, prints it with
  `--version`, and reports the old and new version when `--force` updates an
  existing install.
- `install.sh --uninstall`: removes the generated rule files, the stack-agents
  scripts and the skills this repository ships, while keeping `AGENTS.md`
  (with the project's own overrides), `.ai/conventions.md` and any custom
  skill. Honours `--dry-run`.
- `sync.sh --diff`: like `--check`, but prints the difference for each stale
  file so drift is visible, not just named.
- `CONTRIBUTING.md` with the workflow and the pull-request checklist.
- Tests for `sync.sh` itself: drift detection, regeneration, the Cursor
  frontmatter and the source banner on every generated file.
- A bilingual parity check over `docs/en` and `docs/vi`: same filenames and
  the same number of headings on both sides.
- CI now runs ShellCheck over the shell scripts and Ruff over
  `scripts/load_skill.py`, before the existing drift check and test suite.
