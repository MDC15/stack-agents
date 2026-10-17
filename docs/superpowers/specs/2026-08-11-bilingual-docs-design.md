# pstack-universal: bilingual docs & repo completeness — Design

Date: 2026-08-11
Status: approved in chat (layout + scope), license and git pending

## Intent

A stranger should be able to pick up this repo, understand it, and install it
into a project without asking anything — in English or Vietnamese. Today the
README is Vietnamese-only, the agent rules are English-only, and the repo
ships no license, no CI, no tests and no task runner, so nothing stops the
generated agent files from drifting away from `AGENTS.md`.

## Goals

- README and docs available as two clearly linked versions, English and Vietnamese.
- One command installs into a project (done: `scripts/install.sh`).
- One command checks everything: `bash scripts/test.sh` + `bash scripts/sync.sh --check`.
- CI runs that check so generated agent files cannot drift.
- Repo ships a license, a dependency-free test script, and a Makefile.

## Non-goals

- Translating `AGENTS.md` or anything under `.ai/`. Those are loaded into the
  model's context every session; a bilingual copy doubles token cost and
  contradicts the repo's own "keep AGENTS.md short, load on demand" rule.
- Adding a git remote, or flattening the nested `pstack-universal/pstack-universal/`
  directory. Both are the owner's call, not mine.
- Any new runtime dependency (principle 10).

## Constraints

- No new dependencies. Tests are a plain bash script.
- English is canonical for docs; the Vietnamese files are translations of them.
- Filenames under `docs/en/` and `docs/vi/` mirror each other exactly.
- `bash scripts/sync.sh --check` must stay green at the end of every task.
- Generated files (`CLAUDE.md`, `.clinerules`, `.roorules`, `.windsurfrules`,
  `.github/copilot-instructions.md`, `.cursor/rules/pstack.mdc`) are never
  edited by hand.

## Decisions

1. **English is the source of truth for documentation.** `README.md` is
   English, `README.vi.md` is the Vietnamese translation.
2. **`docs/en/` and `docs/vi/` mirror filenames**, so a reader can swap one
   path segment to switch language.
3. **Agent-facing files stay English-only** (see Non-goals).
4. **Tests are `scripts/test.sh`**, a bash script with plain asserts, covering
   `install.sh` behaviour, `load_skill.py` CLI, and relative-link integrity in
   Markdown. No pytest, no bats.
5. **CI is a single GitHub Actions workflow** running `scripts/test.sh` and
   `scripts/sync.sh --check`.
6. **Docs split into five topics** in both languages: `getting-started`,
   `how-it-works`, `skills`, `custom-agents`, `extending`.
7. **Per-agent rule-file names carry an explicit caveat** that tool versions
   change these paths; the current README already flags this as unverified.

## Open (only the owner can decide)

- ~~License: MIT vs Apache-2.0~~ — **decided 2026-08-11: MIT**.
- `git init` and the eventual remote; commit steps are blocked until then.
- Whether to flatten the nested directory.

## What is left to the builder

The prose of every document. Source material for each topic is named in the
plan: the existing `README.md`, `AGENTS.md`, `.ai/principles.md`, the eight
skill files, the two playbooks, `scripts/sync.sh`, `scripts/install.sh` and
`scripts/load_skill.py`.
