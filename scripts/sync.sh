#!/usr/bin/env bash
# Generate agent-specific rule files from AGENTS.md (the single source of truth).
# Usage: bash scripts/sync.sh [--check | --diff]
#   --check  exit 1 if any generated file differs from AGENTS.md (for CI)
#   --diff   like --check, but print the difference for each stale file
set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
src="$root/AGENTS.md"
check=0
diffmode=0
for arg in "$@"; do
  case "$arg" in
    --check) check=1 ;;
    --diff) check=1; diffmode=1 ;;
    *) echo "unknown option: $arg" >&2; exit 2 ;;
  esac
done

[[ -f "$src" ]] || { echo "missing $src" >&2; exit 1; }

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

banner='<!-- Generated from AGENTS.md by scripts/sync.sh. Edit AGENTS.md, not this file. -->'

{ printf '%s\n\n' "$banner"; cat "$src"; } > "$tmp/plain.md"
{
  printf -- '---\ndescription: pstack-universal operating rules\nalwaysApply: true\n---\n\n'
  cat "$tmp/plain.md"
} > "$tmp/cursor.mdc"

status=0
place() {
  local dest="$root/$1" from="$2"
  if (( check )); then
    if ! cmp -s "$from" "$dest"; then
      echo "out of date: $1"
      if (( diffmode )); then
        diff -L "$1 (current)" -L "$1 (from AGENTS.md)" "$dest" "$from" || true
      fi
      status=1
    fi
  else
    mkdir -p "$(dirname "$dest")"
    cp "$from" "$dest"
    echo "wrote $1"
  fi
}

place CLAUDE.md                        "$tmp/plain.md"
place .clinerules                      "$tmp/plain.md"
place .roorules                        "$tmp/plain.md"
place .windsurfrules                   "$tmp/plain.md"
place .github/copilot-instructions.md  "$tmp/plain.md"
place .cursor/rules/pstack.mdc         "$tmp/cursor.mdc"

(( check && status == 0 )) && echo "all generated files are up to date"
exit "$status"
