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
  link_files=("$root/README.md" "$root/README.vi.md")
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

# --- report ----------------------------------------------------------------
echo "----"
echo "passed: $passed  failed: $failed"
((failed == 0)) || exit 1
