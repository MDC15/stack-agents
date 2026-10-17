"""Framework-agnostic loader for custom agents (LangChain, CrewAI, your own loop).

Wire it up in two steps:
  1. Use system_prompt() as the agent's system prompt.
  2. Expose load_skill(name) as a tool so the agent reads a skill only when needed.

CLI check:  python3 -I scripts/load_skill.py list
            python3 -I scripts/load_skill.py skills tdd
"""
import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
AI_DIR = ROOT / ".ai"
KINDS = ("skills", "playbooks")
NAME_RE = re.compile(r"^[a-z0-9][a-z0-9-]*$")


def system_prompt() -> str:
    return (ROOT / "AGENTS.md").read_text(encoding="utf-8")


def list_modules() -> dict[str, list[str]]:
    return {
        kind: sorted(p.stem for p in (AI_DIR / kind).glob("*.md"))
        for kind in KINDS
    }


def load_skill(name: str, kind: str = "skills") -> str:
    """Return the text of .ai/<kind>/<name>.md. Rejects anything but plain names."""
    if kind not in KINDS:
        raise ValueError(f"kind must be one of {KINDS}")
    if not NAME_RE.fullmatch(name):
        raise ValueError(f"invalid name: {name!r}")
    path = AI_DIR / kind / f"{name}.md"
    if not path.is_file():
        raise FileNotFoundError(f"no such {kind[:-1]}: {name}")
    return path.read_text(encoding="utf-8")


if __name__ == "__main__":
    if len(sys.argv) == 2 and sys.argv[1] == "list":
        for kind, names in list_modules().items():
            print(f"{kind}: {', '.join(names)}")
    elif len(sys.argv) == 3:
        print(load_skill(sys.argv[2], sys.argv[1]))
    else:
        sys.exit(__doc__)
