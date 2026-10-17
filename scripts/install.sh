#!/usr/bin/env bash
# Install stack-agents into a project so nothing has to be copied by hand.
# Usage: bash scripts/install.sh [options] <target-dir>
#
# Options:
#   --force      merge the current AGENTS.md, .ai/, scripts/ and VERSION into the target
#   --dry-run    print what would happen, change nothing
#   --no-sync    install files but do not generate agent config files
#   --uninstall  remove stack-agents-managed files from the target, keep everything else
#   --version    print the version of this checkout and exit
#   -h, --help   show this help
#
# Without --force, existing files are left untouched. On a merge, files the source
# does not contain (such as .ai/conventions.md and custom skills) are kept.
set -euo pipefail

source_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
items=(AGENTS.md .ai scripts VERSION)

force=0
dry=0
run_sync=1
uninstall_mode=0
target=""

usage() {
  # Print this file's leading comment block, up to the first non-comment line.
  awk 'NR > 1 && !/^#/ { exit } NR > 1 { sub(/^# ?/, ""); print }' "${BASH_SOURCE[0]}"
}

# Remove only what stack-agents installed: generated rule files still carrying the
# banner, the stack-agents scripts, the files this source checkout ships under .ai/
# and the skills it ships, and VERSION. The project's own files -- AGENTS.md
# with its Project overrides, .ai/conventions.md, custom skills -- are kept.
uninstall_target() {
  local dir="$1" removed=0 item name src_file
  local banner='Generated from AGENTS.md by scripts/sync.sh'
  # The Cursor file carries frontmatter before its banner, hence head -n 8.
  local generated=(CLAUDE.md .clinerules .roorules .windsurfrules .github/copilot-instructions.md .cursor/rules/stack-agents.mdc)
  local shipped=(sync.sh install.sh test.sh load_skill.py)

  [[ -d "$dir" ]] || { echo "not a directory: $dir" >&2; return 1; }

  for item in "${generated[@]}"; do
    if [[ -f "$dir/$item" ]] && head -n 8 "$dir/$item" | grep -qF "$banner"; then
      ((dry)) || rm "$dir/$item"
      echo "removed $item"
      removed=$((removed + 1))
    fi
  done

  for item in "${shipped[@]}"; do
    if [[ -f "$dir/scripts/$item" ]]; then
      ((dry)) || rm "$dir/scripts/$item"
      echo "removed scripts/$item"
      removed=$((removed + 1))
    fi
  done

  # Top-level files this repo ships under .ai/ (principles.md). The project's
  # own .ai/conventions.md is not in the source and stays.
  for src_file in "$source_root"/.ai/*.md; do
    name="$(basename "$src_file")"
    [[ -f "$dir/.ai/$name" ]] || continue
    ((dry)) || rm "$dir/.ai/$name"
    echo "removed .ai/$name"
    removed=$((removed + 1))
  done

  # Only the skills and playbooks that came from this source checkout; anything
  # else in those directories is the project's own and stays.
  for item in skills playbooks; do
    [[ -d "$dir/.ai/$item" ]] || continue
    for src_file in "$source_root/.ai/$item/"*.md; do
      name="$(basename "$src_file")"
      [[ -f "$dir/.ai/$item/$name" ]] || continue
      ((dry)) || rm "$dir/.ai/$item/$name"
      echo "removed .ai/$item/$name"
      removed=$((removed + 1))
    done
  done

  if [[ -f "$dir/VERSION" ]]; then
    ((dry)) || rm "$dir/VERSION"
    echo "removed VERSION"
    removed=$((removed + 1))
  fi

  if (( ! dry )); then
    rmdir "$dir/.ai/skills" "$dir/.ai/playbooks" 2>/dev/null || true
    rmdir "$dir/.cursor/rules" 2>/dev/null || true
    rmdir "$dir/scripts" "$dir/.cursor" "$dir/.github" 2>/dev/null || true
    rmdir "$dir/.ai" 2>/dev/null || true
  fi

  if ((removed == 0)); then
    echo "nothing to uninstall in $dir"
    return 0
  fi
  if ((dry)); then
    echo "would remove $removed stack-agents-managed files from $dir"
  else
    echo "removed $removed stack-agents-managed files from $dir"
  fi
  cat <<EOF
Kept on purpose: AGENTS.md (it holds this project's Project overrides),
.ai/conventions.md, and any file stack-agents did not ship. Remove them by hand if
you no longer want them.
EOF
}

while (($#)); do
  case "$1" in
    --force) force=1 ;;
    --dry-run) dry=1 ;;
    --no-sync) run_sync=0 ;;
    --uninstall) uninstall_mode=1 ;;
    --version) echo "stack-agents $(cat "$source_root/VERSION")"; exit 0 ;;
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

if ((uninstall_mode)); then
  uninstall_target "$target"
  exit 0
fi

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
old_version=""
[[ -f "$target/VERSION" ]] && old_version="$(cat "$target/VERSION")"
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

new_version="$(cat "$source_root/VERSION")"
if [[ -n "$old_version" && "$old_version" != "$new_version" ]]; then
  echo "updated stack-agents $old_version -> $new_version"
fi

cat <<EOF

stack-agents v$new_version installed in $target
Next: fill in the build, test and lint commands under "Project overrides" in
$target/AGENTS.md, then re-run sync.sh if you edit AGENTS.md.
EOF
