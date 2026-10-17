#!/usr/bin/env bash
# Install pstack-universal into a project so nothing has to be copied by hand.
# Usage: bash scripts/install.sh [options] <target-dir>
#
# Options:
#   --force     merge the current AGENTS.md, .ai/ and scripts/ into the target
#   --dry-run   print what would happen, change nothing
#   --no-sync   install files but do not generate agent config files
#   -h, --help  show this help
#
# Without --force, existing files are left untouched. On a merge, files the source
# does not contain (such as .ai/conventions.md) are kept.
set -euo pipefail

source_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
items=(AGENTS.md .ai scripts)

force=0
dry=0
run_sync=1
target=""

usage() {
  sed -n '2,12p' "${BASH_SOURCE[0]}" | sed 's/^# \{0,1\}//'
}

while (($#)); do
  case "$1" in
    --force) force=1 ;;
    --dry-run) dry=1 ;;
    --no-sync) run_sync=0 ;;
    -h|--help) usage; exit 0 ;;
    -*) echo "unknown option: $1" >&2; usage >&2; exit 2 ;;
    *) target="$1" ;;
  esac
  shift
done

if [[ -z "$target" ]]; then
  echo "missing <target-dir>" >&2
  usage >&2
  exit 2
fi

# Resolve the target even if it does not exist yet.
if [[ -d "$target" ]]; then
  target="$(cd "$target" && pwd)"
else
  parent="$(cd "$(dirname "$target")" && pwd)"
  target="$parent/$(basename "$target")"
fi

[[ -f "$source_root/AGENTS.md" ]] || { echo "missing $source_root/AGENTS.md" >&2; exit 1; }
[[ "$target" != "$source_root" ]] || { echo "target is the source repo; pick another directory" >&2; exit 1; }

conflicts=()
for item in "${items[@]}"; do
  [[ -e "$target/$item" ]] && conflicts+=("$item")
done
if ((${#conflicts[@]})) && ((!force)); then
  echo "already present in $target: ${conflicts[*]}" >&2
  echo "re-run with --force to overwrite, or install into an empty directory" >&2
  exit 1
fi

if ((dry)); then
  echo "would create $target"
  for item in "${items[@]}"; do
    echo "would copy $item"
  done
  ((run_sync)) && echo "would run scripts/sync.sh in $target"
  exit 0
fi

mkdir -p "$target"
for item in "${items[@]}"; do
  if [[ -d "$source_root/$item" ]]; then
    # Merge into the existing directory. Copying the directory itself would
    # nest it (target/.ai/.ai) and leave the real contents untouched on update.
    mkdir -p "$target/$item"
    cp -R "$source_root/$item/." "$target/$item/"
  else
    cp "$source_root/$item" "$target/$item"
  fi
  echo "copied $item"
done

if ((run_sync)); then
  bash "$target/scripts/sync.sh"
fi

cat <<EOF

pstack-universal installed in $target
Next: fill in the build, test and lint commands under "Project overrides" in
$target/AGENTS.md, then re-run sync.sh if you edit AGENTS.md.
EOF
