# pstack-universal: bilingual docs & repo completeness Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Ship the repo as a self-checking, documented, English-and-Vietnamese project: one command installs it, one command verifies it, CI keeps generated agent files from drifting.

**Architecture:** `AGENTS.md` stays the single English source of truth; `scripts/sync.sh` derives the per-agent rule files from it. Documentation becomes a mirrored pair — `docs/en/` canonical, `docs/vi/` its translation — linked from `README.md` (English) and `README.vi.md` (Vietnamese). Verification is a dependency-free `scripts/test.sh` that exercises `install.sh`, the `load_skill.py` CLI, and every relative Markdown link.

**Tech Stack:** bash (no new dependencies), Python 3 stdlib only, GitHub Actions, GNU Make.

**Spec:** `docs/superpowers/specs/2026-08-11-bilingual-docs-design.md`

## Global Constraints

- No new dependencies of any kind (spec: principle 10). Tests are plain bash.
- English is canonical for docs; every `docs/vi/` file is a translation of the `docs/en/` file with the same name.
- Filenames under `docs/en/` and `docs/vi/` mirror each other exactly: `getting-started.md`, `how-it-works.md`, `skills.md`, `custom-agents.md`, `extending.md`.
- `AGENTS.md` and everything under `.ai/` remain English-only and are not translated.
- Generated files are never hand-edited: `CLAUDE.md`, `.clinerules`, `.roorules`, `.windsurfrules`, `.github/copilot-instructions.md`, `.cursor/rules/pstack.mdc`. Only `sync.sh` writes them.
- `bash scripts/sync.sh --check` must exit 0 at the end of every task.
- No git repository exists yet, so commit steps are omitted; they resume once `git init` is agreed with the owner.
- Target toolchain: `bash` 4+, `python3`, `make`.

## Review Focus

1. **Relative links between the two languages and into `docs/` break silently** (README.md → README.vi.md → docs/vi/…), leaving a reader on a 404. Expected: every relative link resolves to an existing file. Test owned by Task 9: parse `[text](target)` in all Markdown, skip `http(s)://`, `#`, `mailto:`, strip anchors, assert the file exists.
2. **Editing `AGENTS.md` leaves the generated agent files stale**, so Claude Code/Cursor/Cline silently run old rules. Expected: mismatch is detected. Test owned by Task 2, re-run by Task 3/4: `bash scripts/sync.sh --check` exits 0.
3. **`install.sh --force` deletes a project's learned `.ai/conventions.md`** on update, destroying user data. Expected: it survives. Test owned by Task 1: write a marker file, `install --force`, assert content unchanged.
4. **`load_skill.py` can be handed a crafted name to read a file outside `.ai/`** (`../../etc/passwd`). Expected: non-zero exit, no file content printed. Test owned by Task 1.
5. **`make install` with no `TARGET` does something surprising** instead of refusing. Expected: exit 2 and a usage line. Test owned by Task 3.

---

### Task 1: Dependency-free test script

**Files:**
- Create: `scripts/test.sh`

**Interfaces:**
- Consumes: `scripts/install.sh` CLI contract — exit `0` success, `1` refusal (existing target / target is source), `2` usage error (missing target / unknown flag); flags `--force`, `--dry-run`, `--no-sync`. Consumes `python3 -I scripts/load_skill.py` CLI — `list` prints `skills: …` and `playbooks: …`; `<kind> <name>` prints that file; any invalid kind or name exits non-zero.
- Produces: `bash scripts/test.sh` → prints one `PASS:`/`FAIL:` line per case, exits `0` when all pass and `1` otherwise. Consumed by Makefile `check` (Task 3) and CI (Task 4).

- [ ] **Step 1: Write `scripts/test.sh`**

Shebang `#!/usr/bin/env bash`, `set -euo pipefail`. Helpers: `pass()`/`fail()` incrementing counters, `tmp="$(mktemp -d)"` with `trap 'rm -rf "$tmp"' EXIT`, and `expect_exit <want> <cmd…>` asserting an exit code. `src` = the repo root resolved from `BASH_SOURCE`. Cases, in this order:

```
bash -n scripts/install.sh
bash -n scripts/sync.sh
fresh install      : install.sh "$tmp/p1" exits 0; $tmp/p1/AGENTS.md exists;
                     (cd $tmp/p1 && bash scripts/sync.sh --check) exits 0
refuse without flag: install.sh "$tmp/p1" exits 1; stderr contains "already present"
force preserves    : echo "# project rule" > $tmp/p1/.ai/conventions.md;
                     install.sh --force "$tmp/p1" exits 0;
                     $tmp/p1/.ai/conventions.md still contains "# project rule"
dry-run creates    : install.sh --dry-run "$tmp/p2" exits 0; $tmp/p2 does not exist
no-sync             : install.sh --no-sync "$tmp/p3" exits 0;
                     $tmp/p3/AGENTS.md exists; $tmp/p3/CLAUDE.md does not
target==source      : install.sh "$src" exits 1
no args             : install.sh exits 2
unknown flag        : install.sh --nope "$tmp/x" exits 2
loader list         : python3 -I scripts/load_skill.py list → stdout contains
                      "skills: architect, arena, checkpoint, conventions, red-team, tdd, unslop, verify"
loader one skill    : python3 -I scripts/load_skill.py skills tdd → stdout contains "/tdd"
loader playbook     : python3 -I scripts/load_skill.py playbooks debug → stdout contains "Playbook: debug"
loader bad name     : python3 -I scripts/load_skill.py skills ../../etc/passwd exits non-zero
loader bad kind     : python3 -I scripts/load_skill.py bogus tdd exits non-zero
```

The `expect_exit` helper must capture stderr to a file so refusal messages can be asserted without `set -e` aborting on non-zero.

- [ ] **Step 2: Run the script and confirm every case passes**

Run: `bash scripts/test.sh`
Expected: one `PASS:` line per case above, final line reporting the total, exit `0`.

- [ ] **Step 3: Make it executable and confirm the failure path works**

Run: `chmod +x scripts/test.sh && bash scripts/test.sh`
Expected: still exit `0`. Then confirm the harness actually detects failure: temporarily change one expected exit code in `expect_exit`, re-run, expect exit `1` with a `FAIL:` line, and revert.

---

### Task 2: Fill in this repo's own Project overrides and re-sync

**Files:**
- Modify: `AGENTS.md` (the four `<fill in>` lines under `## Project overrides`)
- Generated: all six rule files, via `scripts/sync.sh`

**Interfaces:**
- Consumes: `AGENTS.md` is the sole input to `sync.sh`.
- Produces: `AGENTS.md` documents this repo's real commands — `Test: bash scripts/test.sh`, `Lint/typecheck: bash scripts/sync.sh --check` — which Task 5 and Task 7 quote in `getting-started.md`.

- [ ] **Step 1: Replace the placeholders**

Under `## Project overrides`, set:
- `Build: none (this repo is distributed as Markdown and scripts)`
- `Test: bash scripts/test.sh`
- `Lint/typecheck: bash scripts/sync.sh --check`

- [ ] **Step 2: Regenerate the agent files**

Run: `bash scripts/sync.sh`
Expected: six `wrote …` lines (`CLAUDE.md`, `.clinerules`, `.roorules`, `.windsurfrules`, `.github/copilot-instructions.md`, `.cursor/rules/pstack.mdc`).

- [ ] **Step 3: Confirm no drift remains**

Run: `bash scripts/sync.sh --check`
Expected: `all generated files are up to date`, exit `0`.

---

### Task 3: Makefile

**Files:**
- Create: `Makefile`

**Interfaces:**
- Consumes: `bash scripts/test.sh` (Task 1), `bash scripts/sync.sh [--check]`, `bash scripts/install.sh`.
- Produces: `make check` → runs `sync.sh --check` then `scripts/test.sh`; `make sync`; `make install TARGET=<dir>`; `make help` (default target, lists them).

- [ ] **Step 1: Write `Makefile`**

`.PHONY: help check sync install`. `help` is the first target (so bare `make` prints it). `install` must fail with exit `2` and the line `usage: make install TARGET=/path/to/project` when `TARGET` is empty.

- [ ] **Step 2: Verify `check` is green**

Run: `make check`
Expected: `all generated files are up to date`, then all test `PASS:` lines, exit `0`.

- [ ] **Step 3: Verify the guard on `install`**

Run: `make install; echo $?`
Expected: prints the usage line, exit `2`.

- [ ] **Step 4: Verify `help` and `sync`**

Run: `make help` and `make sync`
Expected: `help` lists the targets; `sync` writes six files and `make check` still passes afterwards.

---

### Task 4: CI workflow

**Files:**
- Create: `.github/workflows/check.yml`

**Interfaces:**
- Consumes: `make check` from Task 3.
- Produces: a workflow named `check` that runs on `push` and `pull_request`.

- [ ] **Step 1: Write the workflow**

`actions/checkout@v4`, then a single `run: make check` step. `runs-on: ubuntu-latest`.

- [ ] **Step 2: Validate the YAML parses**

Run: `python3 -c 'import yaml,sys; yaml.safe_load(open(".github/workflows/check.yml")); print("yaml ok")'`
Expected: `yaml ok`. If PyYAML is unavailable locally, say so and fall back to a careful read — do not claim it validated.

- [ ] **Step 3: Confirm the workflow's job passes locally by running its exact command**

Run: `make check`
Expected: exit `0`.

---

### Task 5: LICENSE

**Files:**
- Create: `LICENSE`
- Modify: `README.md` (License section, if it links to the file)

**Interfaces:**
- Consumes: the license choice recorded in the spec's *Open* section (MIT or Apache-2.0).
- Produces: a `LICENSE` file whose text is the verbatim standard text for the chosen SPDX id, with the copyright line `Copyright (c) 2026 pstack-universal contributors`.

- [ ] **Step 1: Record the owner's choice in the spec**

Update the *Open* section of `docs/superpowers/specs/2026-08-11-bilingual-docs-design.md` to mark the license decision as answered, naming the SPDX id.

- [ ] **Step 2: Write `LICENSE`**

Use the verbatim standard text for the chosen license. Do not paraphrase it.

- [ ] **Step 3: Verify it matches the standard text**

Run: `head -3 LICENSE` and `wc -l LICENSE`
Expected: the opening line matches the chosen license name exactly, and the line count matches the standard text (MIT is 21 lines).

---

### Task 6: English README (canonical)

**Files:**
- Create: `README.md` (replacing the current Vietnamese one; preserve nothing — its content becomes the basis of `README.vi.md` in Task 7)

**Interfaces:**
- Consumes: the current `README.md` content, `AGENTS.md`, `scripts/install.sh`, `scripts/sync.sh`.
- Produces: the canonical English README, and the exact section list `README.vi.md` must mirror.

- [ ] **Step 1: Rewrite `README.md` in English**

Sections, in this order:
1. `# pstack-universal` plus one-paragraph description.
2. `**Languages:** English | [Tiếng Việt](README.vi.md)` — one line, directly under the title.
3. `## How it works` — one source of truth, `sync.sh` generates per-tool files, skills load on demand.
4. `## Structure` — tree including `scripts/install.sh`, `scripts/test.sh`, `Makefile`, `docs/en/`, `docs/vi/`.
5. `## Install` — the one-line `install.sh` command, the options table (`--force`, `--dry-run`, `--no-sync`), the overwrite-refusal note, the `conventions.md` preservation note, and a `<details>` block with the manual steps.
6. `## Generated files` — table mapping each agent to its file, with the caveat that tool versions change these paths.
7. `## Commands` — the eight `/…` commands plus `/auto`, one line each.
8. `## Custom agents` — the `load_skill` snippet, linking to `docs/en/custom-agents.md`.
9. `## Extending` — one line, linking to `docs/en/extending.md`.
10. `## Verification` — `make check`, what it runs.
11. `## Limitations` — keep the existing two honest bullets.
12. `## License`.

Target length: 100–120 lines. Link every topic to its `docs/en/` page.

- [ ] **Step 2: Confirm nothing was lost from the original**

Run: `git diff README.md` if git exists, otherwise diff against a copy taken in Step 1.
Expected: every fact from the original (agent→file table, command table, limitations, install options) appears; no Vietnamese text remains outside the language-switch line.

---

### Task 7: Vietnamese README

**Files:**
- Create: `README.vi.md`

**Interfaces:**
- Consumes: `README.md` from Task 6 (section-for-section translation).
- Produces: the Vietnamese mirror; both files carry a cross-link line.

- [ ] **Step 1: Translate `README.md` into `README.vi.md`**

Same section list, same tables, same code blocks verbatim (commands and paths are never translated).

- [ ] **Step 2: Set the language-switch line to point back**

In `README.vi.md`, line under the title: `**Ngôn ngữ:** [English](README.md) | Tiếng Việt`.
In `README.md`, keep `English | [Tiếng Việt](README.vi.md)`.

- [ ] **Step 3: Confirm both cross-links resolve**

Run: `test -f README.md && test -f README.vi.md && grep -q 'README.vi.md' README.md && grep -q 'README.md' README.vi.md && echo ok`
Expected: `ok`.

---

### Task 8: English docs

**Files:**
- Create: `docs/en/getting-started.md`
- Create: `docs/en/how-it-works.md`
- Create: `docs/en/skills.md`
- Create: `docs/en/custom-agents.md`
- Create: `docs/en/extending.md`

**Interfaces:**
- Consumes: `README.md`, `AGENTS.md`, `scripts/install.sh`, `scripts/sync.sh`, `scripts/load_skill.py`, `.ai/skills/*.md`, `.ai/playbooks/*.md`, `.ai/principles.md`.
- Produces: the five canonical English topics. Task 9 translates these exact filenames.

- [ ] **Step 1: `docs/en/getting-started.md`** (40–70 lines)

Sections: Prerequisites (`bash`, `python3`) → Install (script, options table, manual fallback) → Fill in *Project overrides* → Regenerate (`bash scripts/sync.sh`) → Verify (`make check`, or `bash scripts/test.sh`) → What your agent now sees → Next: `how-it-works.md`.

- [ ] **Step 2: `docs/en/how-it-works.md`** (40–70 lines)

Sections: Why one source of truth → What `sync.sh` generates (table, including that the Cursor file adds `alwaysApply: true` frontmatter) → Why skills load on demand instead of being inlined → What `AGENTS.md` is and is not (it instructs, it does not enforce) → The `install.sh`/`sync.sh` split (install copies, sync derives).

- [ ] **Step 3: `docs/en/skills.md`** (80–110 lines)

One `### /command` section per skill, in the order they appear in `AGENTS.md`: `architect`, `arena`, `red-team` (note the `/swarm` alias), `tdd`, `verify`, `checkpoint`, `unslop`, `conventions`. Then `### /auto <goal>` covering the two playbooks `debug` and `refactor`. Each section: when to trigger it (from the `AGENTS.md` table), what it does, and what it outputs. Keep it a summary — the `.ai/` files remain the authority.

- [ ] **Step 4: `docs/en/custom-agents.md`** (40–60 lines)

Sections: Why (LangChain/CrewAI/custom loop) → `system_prompt()` → `list_modules()` → `load_skill(name, kind)` → the name-safety guarantee (regex `^[a-z0-9][a-z0-9-]*$`, `kind` limited to `skills`/`playbooks`, so no path escape) → CLI checks (`python3 -I scripts/load_skill.py list`, `… skills tdd`) → Wiring it as a tool.

- [ ] **Step 5: `docs/en/extending.md`** (40–70 lines)

Sections: Add a skill (create `.ai/skills/<name>.md`, add a row to the `AGENTS.md` table, run `sync.sh`) → Add a playbook → Keep `AGENTS.md` short → Verify (`make check`) → Contribute (test script, CI expectation).

- [ ] **Step 6: Confirm each file covers its source material**

Run: `for f in docs/en/*.md; do echo "== $f"; grep -c '^#' "$f"; wc -l < "$f"; done`
Expected: five files, each with a title heading and at least two `##` sections, and each within its line budget above.

---

### Task 9: Vietnamese docs

**Files:**
- Create: `docs/vi/getting-started.md`
- Create: `docs/vi/how-it-works.md`
- Create: `docs/vi/skills.md`
- Create: `docs/vi/custom-agents.md`
- Create: `docs/vi/extending.md`

**Interfaces:**
- Consumes: the five `docs/en/` files from Task 8.
- Produces: a mirrored `docs/vi/` tree; every file carries an English link and every `docs/en/` file carries a Vietnamese link.

- [ ] **Step 1: Translate each `docs/en/` file into `docs/vi/` with the identical filename**

Same sections, same tables, code blocks and commands verbatim.

- [ ] **Step 2: Add cross-links at the top of every docs file**

In `docs/en/<name>.md`: `> Tiếng Việt: [docs/vi/<name>.md](../../docs/vi/<name>.md)`.
In `docs/vi/<name>.md`: `> English: [docs/en/<name>.md](../../docs/en/<name>.md)`.

- [ ] **Step 3: Confirm the trees mirror**

Run: `diff <(cd docs/en && ls) <(cd docs/vi && ls) && echo mirror ok`
Expected: `mirror ok`.

---

### Task 10: Link-integrity check and final verification

**Files:**
- Modify: `scripts/test.sh` (append the link-check case)

**Interfaces:**
- Consumes: every Markdown file written in Tasks 6–9.
- Produces: a `test.sh` that also fails when any relative Markdown link dangles, closing Review Focus item 1.

- [ ] **Step 1: Append a link-integrity case to `scripts/test.sh`**

Walk exactly `README.md`, `README.vi.md`, `docs/en/*.md` and `docs/vi/*.md` — not `docs/superpowers/**`, which is internal planning material. Strip fenced code blocks (``` … ```) and inline spans in backticks first, so a literal `` `[text](target)` `` is never treated as a link. Then, for each remaining `[text](target)`, skip targets starting with `http://`, `https://`, `mailto:` or `#`, strip any `#anchor`, and assert `[[ -e "$dir/$target" ]]` where `$dir` is the directory of the containing file. On failure print the file, the line and the dangling target.

- [ ] **Step 2: Run it and fix any dangling link it finds**

Run: `bash scripts/test.sh`
Expected: all cases pass, exit `0`. If a link dangles, correct it in the Markdown, not in the checker.

- [ ] **Step 3: Confirm the checker really fails on a bad link**

Run: create `docs/en/zz-probe.md` containing `[x](does-not-exist.md)`, run `bash scripts/test.sh`, then delete the probe.
Expected: exit `1` naming `docs/en/zz-probe.md`; after deletion, exit `0`.

- [ ] **Step 4: Run the full gate**

Run: `make check`
Expected: `all generated files are up to date`, every `PASS:` line, exit `0`.
