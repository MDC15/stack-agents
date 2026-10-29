# Custom agents

> Tiếng Việt: [docs/vi/custom-agents.md](../../docs/vi/custom-agents.md)

The generated rule files cover the tools that read a fixed path. If you run your
own loop — LangChain, CrewAI, or something you wrote — `scripts/load_skill.py`
gives you the same two pieces: the rules as a system prompt, and the skills as a
tool the model calls when it needs one.

## `system_prompt()`

```python
from scripts.load_skill import system_prompt

prompt = system_prompt()   # the text of AGENTS.md
```

Returns the contents of `AGENTS.md`. Use it as the agent's system prompt.

## `list_modules()`

```python
from scripts.load_skill import list_modules

list_modules()
# {'skills': ['architect', 'arena', ...], 'playbooks': ['debug', 'refactor']}
```

Returns the available names, sorted, grouped by kind. Useful for building the
tool description or a help command.

## `load_skill(name, kind="skills")`

```python
from scripts.load_skill import load_skill

load_skill("tdd")                     # text of .ai/skills/tdd.md
load_skill("debug", kind="playbooks") # text of .ai/playbooks/debug.md
```

`kind` is `"skills"` (the default) or `"playbooks"`.

## The name check is the safety property

`load_skill` accepts only names matching `^[a-z0-9][a-z0-9-]*$`, and only the two
kinds above. A name like `../../etc/passwd` raises `ValueError`, so a model that
emits a hostile or garbled tool argument cannot use this function to read a file
outside `.ai/`.

This is the reason the restriction exists rather than a cosmetic regex. If you
re-implement the loader, keep it.

## Check it from the shell

```bash
python3 -I scripts/load_skill.py list
python3 -I scripts/load_skill.py skills tdd
python3 -I scripts/load_skill.py playbooks debug
```

`-I` isolates the interpreter from the environment, which is what you want for a
sanity check. Passing an unknown kind or an invalid name exits non-zero.

## Wiring it as a tool

Expose one tool whose schema is `{name: string, kind: "skills" | "playbooks"}`
and have its body call `load_skill`, returning the text. Leave the description
telling the model to call it only when a command is invoked, so the behaviour
matches the operating loop in `AGENTS.md`.

The repository's own test suite covers the name check — it is
[`repo_test.go`](../../repo_test.go) now, run by `go test ./...`; and
[cli.md](cli.md) documents the same CLI the tests drive.

## Next

- [how-it-works.md](how-it-works.md) — the sync model behind the tool files.
- [skills.md](skills.md) — what the loaded files actually contain.