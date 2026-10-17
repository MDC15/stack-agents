#!/usr/bin/env bash
# Self-check for pstack-universal.
#   bash scripts/test.sh     exit 0 = every case passed
# Covers: install.sh behaviour, load_skill.py CLI, Markdown link integrity.
set -uo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

passed=0
failed=0
rc=0

ok() { passed=$((passed + 1)); echo "PASS: $1"; }
no() { failed=$((failed + 1)); echo "FAIL: $1"; }

# run <cmd...>: capture exit code in $rc, stdout in $tmp/out, stderr in $tmp/err
run() {
  rc=0
  "$@" >"$tmp/out" 2>"$tmp/err" || rc=$?
}

want_rc() { # <want> <description>
  if ((rc == $1)); then ok "$2"; else no "$2 (want exit $1, got $rc)"; fi
}
want_out() { # <needle> <description>
  if grep -qF -- "$1" "$tmp/out"; then ok "$2"; else no "$2 (stdout lacks: $1)"; fi
}
want_err() { # <needle> <description>
  if grep -qF -- "$1" "$tmp/err"; then ok "$2"; else no "$2 (stderr lacks: $1)"; fi
}

# --- install.sh ------------------------------------------------------------
run bash -n "$root/scripts/install.sh"
want_rc 0 "syntax: install.sh parses"

run bash -n "$root/scripts/sync.sh"
want_rc 0 "syntax: sync.sh parses"

run "$root/scripts/install.sh" "$tmp/p1"
want_rc 0 "install: fresh target succeeds"
want_rc 0 "install: wrote AGENTS.md" ; [[ -f "$tmp/p1/AGENTS.md" ]] || no "install: wrote AGENTS.md (missing)"

run bash -c "cd '$tmp/p1' && bash scripts/sync.sh --check"
want_rc 0 "install: generated files are in sync inside the target"
want_out "all generated files are up to date" "install: sync --check reports no drift"

run "$root/scripts/install.sh" "$tmp/p1"
want_rc 1 "install: refuses an existing target without --force"
want_err "already present" "install: refusal explains what is already present"

echo "# project rule" >"$tmp/p1/.ai/conventions.md"
run "$root/scripts/install.sh" --force "$tmp/p1"
want_rc 0 "install: --force succeeds on an existing target"
if grep -q "# project rule" "$tmp/p1/.ai/conventions.md"; then
  ok "install: --force preserves .ai/conventions.md"
else
  no "install: --force preserves .ai/conventions.md"
fi

# --force must update directory contents in place, not copy a directory into itself
cp -R "$root" "$tmp/src2"
printf '%s\n' '<!-- updated by the regression test -->' >>"$tmp/src2/.ai/skills/tdd.md"
run "$tmp/src2/scripts/install.sh" --force "$tmp/p1"
want_rc 0 "install: --force update run succeeds"
if grep -q 'updated by the regression test' "$tmp/p1/.ai/skills/tdd.md"; then
  ok "install: --force propagates updated skill content"
else
  no "install: --force propagates updated skill content (target still has the old file)"
fi
if [[ ! -e "$tmp/p1/.ai/.ai" && ! -e "$tmp/p1/scripts/scripts" ]]; then
  ok "install: --force does not nest .ai/ or scripts/ inside themselves"
else
  no "install: --force does not nest .ai/ or scripts/ inside themselves"
fi

run "$root/scripts/install.sh" --dry-run "$tmp/p2"
want_rc 0 "install: --dry-run succeeds"
if [[ ! -e "$tmp/p2" ]]; then ok "install: --dry-run creates nothing"; else no "install: --dry-run creates nothing"; fi

run "$root/scripts/install.sh" --no-sync "$tmp/p3"
want_rc 0 "install: --no-sync succeeds"
if [[ -f "$tmp/p3/AGENTS.md" ]]; then ok "install: --no-sync still copies files"; else no "install: --no-sync still copies files"; fi
if [[ ! -e "$tmp/p3/CLAUDE.md" ]]; then ok "install: --no-sync skips generation"; else no "install: --no-sync skips generation"; fi

run "$root/scripts/install.sh" "$root"
want_rc 1 "install: refuses the source repo as target"

run "$root/scripts/install.sh"
want_rc 2 "install: missing target is a usage error"

run "$root/scripts/install.sh" --nope "$tmp/x"
want_rc 2 "install: unknown flag is a usage error"

run "$root/scripts/install.sh" --version
want_rc 0 "install: --version succeeds"
want_out "pstack-universal $(cat "$root/VERSION")" "install: --version reports this checkout's version"

# --- sync.sh ---------------------------------------------------------------
cp -R "$root" "$tmp/syncsrc"
printf '\n<!-- drift marker -->\n' >>"$tmp/syncsrc/AGENTS.md"

run bash -c "cd '$tmp/syncsrc' && bash scripts/sync.sh --check"
want_rc 1 "sync: --check fails after AGENTS.md is edited"
want_out "out of date: CLAUDE.md" "sync: --check names the stale file"

run bash -c "cd '$tmp/syncsrc' && bash scripts/sync.sh --diff"
want_rc 1 "sync: --diff fails on drift like --check"
want_out "drift marker" "sync: --diff prints the difference"

run bash -c "cd '$tmp/syncsrc' && bash scripts/sync.sh"
want_rc 0 "sync: regenerates after a drift"

run bash -c "cd '$tmp/syncsrc' && bash scripts/sync.sh --check"
want_rc 0 "sync: --check passes again after a sync"

run bash -c "cd '$tmp/syncsrc' && bash scripts/sync.sh --nope"
want_rc 2 "sync: unknown flag is a usage error"

if grep -q '^alwaysApply: true$' "$tmp/syncsrc/.cursor/rules/pstack.mdc"; then
  ok "sync: cursor file keeps the alwaysApply frontmatter"
else
  no "sync: cursor file keeps the alwaysApply frontmatter"
fi

banner_missing=0
for f in CLAUDE.md .clinerules .roorules .windsurfrules .github/copilot-instructions.md .cursor/rules/pstack.mdc; do
  grep -q 'Generated from AGENTS.md' "$tmp/syncsrc/$f" || banner_missing=1
done
if ((banner_missing == 0)); then
  ok "sync: every generated file carries the source banner"
else
  no "sync: every generated file carries the source banner"
fi

# --- uninstall --------------------------------------------------------------
run "$root/scripts/install.sh" "$tmp/p4"
want_rc 0 "uninstall: setup install succeeds"
echo "# project rule" >"$tmp/p4/.ai/conventions.md"
echo "# mine" >"$tmp/p4/.ai/skills/mine.md"

run "$root/scripts/install.sh" --uninstall --dry-run "$tmp/p4"
want_rc 0 "uninstall: --dry-run succeeds"
want_out "would remove" "uninstall: --dry-run says what it would remove"
if [[ -f "$tmp/p4/CLAUDE.md" && -f "$tmp/p4/scripts/sync.sh" ]]; then
  ok "uninstall: --dry-run removes nothing"
else
  no "uninstall: --dry-run removes nothing"
fi

run "$root/scripts/install.sh" --uninstall "$tmp/p4"
want_rc 0 "uninstall: succeeds on an installed target"
want_out "removed CLAUDE.md" "uninstall: reports the generated files it removed"
if [[ ! -e "$tmp/p4/CLAUDE.md" && ! -e "$tmp/p4/.cursor" && ! -e "$tmp/p4/scripts" && ! -e "$tmp/p4/VERSION" && ! -e "$tmp/p4/.ai/principles.md" ]]; then
  ok "uninstall: removes generated files, scripts, VERSION and .ai files"
else
  no "uninstall: removes generated files, scripts, VERSION and .ai files"
fi
if [[ -f "$tmp/p4/AGENTS.md" && -f "$tmp/p4/.ai/conventions.md" && -f "$tmp/p4/.ai/skills/mine.md" ]]; then
  ok "uninstall: keeps AGENTS.md, .ai/conventions.md and custom skills"
else
  no "uninstall: keeps AGENTS.md, .ai/conventions.md and custom skills"
fi

run "$root/scripts/install.sh" --uninstall "$tmp/p4"
want_rc 0 "uninstall: a second run succeeds"
want_out "nothing to uninstall" "uninstall: a second run has nothing to remove"

# --- load_skill.py ---------------------------------------------------------
run python3 -I "$root/scripts/load_skill.py" list
want_rc 0 "loader: list succeeds"
want_out "skills: architect, arena, checkpoint, conventions, red-team, tdd, unslop, verify" "loader: lists every skill"

run python3 -I "$root/scripts/load_skill.py" skills tdd
want_rc 0 "loader: loads a skill by name"
want_out "/tdd" "loader: returns the tdd skill"

run python3 -I "$root/scripts/load_skill.py" playbooks debug
want_rc 0 "loader: loads a playbook by name"
want_out "Playbook: debug" "loader: returns the debug playbook"

run python3 -I "$root/scripts/load_skill.py" skills ../../etc/passwd
if ((rc != 0)); then ok "loader: rejects a path-traversal name"; else no "loader: rejects a path-traversal name (exit 0)"; fi

run python3 -I "$root/scripts/load_skill.py" bogus tdd
if ((rc != 0)); then ok "loader: rejects an unknown kind"; else no "loader: rejects an unknown kind (exit 0)"; fi

# --- Markdown link integrity ------------------------------------------------
# Published documentation only: README pair and both docs/ languages. An installed
# project ships neither, so skip there instead of reporting a false failure.
if [[ -f "$root/README.md" || -d "$root/docs" ]]; then
  link_files=("$root/README.md" "$root/README.vi.md" "$root/CONTRIBUTING.md")
  while IFS= read -r f; do link_files+=("$f"); done < <(find "$root/docs/en" "$root/docs/vi" -name '*.md' 2>/dev/null | sort)

  # Drop fenced blocks and inline code spans so examples are never read as links.
  strip_code() {
    awk '
      /^```/ { fence = !fence; next }
      fence { next }
      { line = $0
        while (match(line, /`[^`]*`/)) line = substr(line, 1, RSTART - 1) substr(line, RSTART + RLENGTH)
        print line }
    ' "$1"
  }

  links_checked=0
  for f in "${link_files[@]}"; do
    [[ -f "$f" ]] || continue
    dir="$(dirname "$f")"
    while IFS= read -r target; do
      [[ -n "$target" ]] || continue
      case "$target" in http://* | https://* | mailto:* | '#'*) continue ;; esac
      target="${target%%#*}"
      [[ -n "$target" ]] || continue
      links_checked=$((links_checked + 1))
      [[ -e "$dir/$target" ]] || no "link: ${f#"$root"/} → $target (missing)"
    done < <(strip_code "$f" | grep -oE '\]\([^)]+\)' | sed 's/^](\(.*\))$/\1/')
  done
  if ((links_checked)); then
    ok "links: $links_checked relative links across the documentation resolve"
  else
    no "links: no relative links found to check"
  fi
else
  echo "SKIP: link check (this checkout has no README.md or docs/)"
fi

# --- bilingual documentation parity -----------------------------------------
# docs/en and docs/vi must carry the same filenames with the same heading
# counts, so a translation cannot silently fall behind. Installed projects
# ship no docs/, so skip there like the link check does.
if [[ -d "$root/docs/en" && -d "$root/docs/vi" ]]; then
  parity_pairs=0
  parity_ok=0
  for f in "$root"/docs/en/*.md; do
    base="$(basename "$f")"
    if [[ ! -f "$root/docs/vi/$base" ]]; then
      no "parity: docs/vi/$base is missing"
      continue
    fi
    parity_pairs=$((parity_pairs + 1))
    en_heads="$(grep -c '^#' "$f")"
    vi_heads="$(grep -c '^#' "$root/docs/vi/$base")"
    if [[ "$en_heads" == "$vi_heads" ]]; then
      parity_ok=$((parity_ok + 1))
    else
      no "parity: $base has $en_heads headings in docs/en but $vi_heads in docs/vi"
    fi
  done
  while IFS= read -r f; do
    base="$(basename "$f")"
    [[ -f "$root/docs/en/$base" ]] || no "parity: docs/en/$base is missing (only docs/vi has it)"
  done < <(find "$root/docs/vi" -maxdepth 1 -name '*.md' | sort)
  if ((parity_pairs && parity_ok == parity_pairs)); then
    ok "parity: $parity_pairs documentation files match across docs/en and docs/vi"
  fi
else
  echo "SKIP: bilingual parity (this checkout has no docs/en or docs/vi)"
fi

# --- report ----------------------------------------------------------------
echo "----"
echo "passed: $passed  failed: $failed"
((failed == 0)) || exit 1
