# Commands

> Tiếng Việt: [docs/vi/skills.md](../../docs/vi/skills.md)

Eight skills and one playbook selector, plus `/auto` to choose a playbook. Each
command maps to a file under `.ai/`; this page summarises them, and the files
remain the authority.

Only `/verify` runs without being asked. The rest run when you invoke them, or
when the agent asks first and you agree.

## /architect

**Use when** the work needs a design before code exists: a new feature, a new
module, a schema or API change.

Restates the goal and constraints in at most three lines, lists assumptions, and
asks only questions that block the design. Reads the code the design touches,
once. Then covers only the sections that apply: components and their
responsibilities, data model (a Mermaid `erDiagram` when there are more than two
entities), data flow, interface contracts, failure modes, migration and rollout,
test strategy. For each significant choice it names the alternative it rejected.

**Output:** one page. If it needs more, the scope is too big and it proposes a
split. It writes no implementation code and waits for approval or changes.

## /arena

**Use when** two or more approaches are viable and the choice is costly to
reverse. Skip it when one option is obviously better.

Picks the criteria that matter for this task — latency, memory, simplicity,
testability, operability, migration cost — rather than a fixed list. Produces
two approaches, three at most, named by what they are ("in-process cache",
"Redis") and differing in a real trade-off. Sketches how each works, its main
risks and what it costs to reverse.

**Output:** one comparison table, one recommendation, and what finding would
change its mind. It implements nothing until you choose. If a candidate is
clearly dominated, it says so instead of padding the table.

## /red-team (alias /swarm)

**Use when** a design or a diff needs adversarial review before it ships.

Attacks through three lenses:

- **Security** — input validation, authn/authz, injection, secrets, data
  exposure, dependency risk.
- **Operations** — failure and retry behaviour, timeouts, load, observability,
  rollback, migrations, config drift.
- **QA** — untested paths, edge cases, concurrency, bad input, backward
  compatibility, ambiguous requirements.

If subagents are available it runs one lens per subagent, so the lenses do not
anchor on each other; otherwise it runs them in sequence and says the
independence is limited.

**Output:** at most five findings per lens, each with severity, evidence (file
and line, or a concrete scenario) and a suggested fix, then a verdict: **ship**,
**fix first**, or **redesign**. It applies no fixes until you pick.

## /tdd

**Use when** the behaviour can be tested and a test harness exists.

One behaviour per cycle: **red** (write one failing test and confirm it fails for
the right reason), **green** (the least code that passes), **refactor** (improve
structure with tests green). Bug fixes start with a regression test.

Never writes production code without a failing test demanding it. Tests behaviour
through public interfaces, not private details. Mocks only at system boundaries —
network, clock, filesystem — and never mocks what it owns unless there is no
cheaper way. Never edits a test to make it pass unless the requirement itself
changed, and says so when it did. If the repo has no harness, it says so and
proposes the smallest one that fits the stack.

## /verify

**Use when** work is about to be reported as done. This one runs automatically.

Writes the acceptance points from the original request as a checklist. Runs the
commands in **Project overrides**, in order: format and lint, typecheck, targeted
tests, the wider suite, build. Then exercises the change directly — starts the
CLI, calls the endpoint, opens the page — because passing tests do not cover UI,
CLI or integration behaviour. Finally reviews its own diff for stray debug
output, unrelated edits, secrets and leftover TODOs.

**Output:** two lists — verified (each command and its result) and not verified
(what, and why). It never says "done", "fixed" or "works" without a command
result behind it.

## /checkpoint

**Use when** a command changes many files or is hard to undo: migrations,
codegen, bulk rewrites, dependency upgrades, formatters run over the whole repo.
Not for read-only commands.

Before: records `git status --short` and `git rev-parse HEAD`, or lists the
relevant directory with modification times outside git, and warns if uncommitted
work could be overwritten. After: compares status and diff again and confirms
only the expected files changed.

**If something unexpected changed** it stops, reports which files and how, and
offers options. It does not run further destructive commands to clean up.

## /unslop

**Use when** prose needs tightening: answers, docs, commit messages, PR
descriptions, comments.

Cuts openers and closers, restating the question, hedging that carries no
information, explanations of what the code already says, and adjectives with no
measurable claim. Keeps names, numbers, versions, commands and paths exactly as
they are, along with caveats that change what the reader will do, decisions and
their reasons, and steps in order.

**Method:** draft, delete everything that does not change what the reader knows
or does, prefer one idea per sentence and active voice, then check no fact or
step was lost. For code it deletes dead code, comments that restate the code and
wrappers that only forward calls — without touching behaviour.

## /conventions

**Use when** a repo's habits should be captured so future sessions follow them
without being told again.

Collects candidate rules from config (linters, formatters, `tsconfig`,
`.editorconfig`, CI workflows), from code (naming, file layout, error handling,
test style across several representative files), from `git log` for commit
style, and from corrections you made in the conversation. Keeps only rules seen
in at least two places or stated by you. Writes each as one imperative line with
its evidence in brackets, groups by topic, and holds the file under 40 lines.

**Output:** proposes `.ai/conventions.md` and saves it only after you approve.
Never includes secrets or personal data, marks unconfirmed rules `(inferred)`,
and merges into an existing file instead of rewriting it.

## /auto <goal>

**Use when** you have a goal but not a command.

Picks the matching playbook under `.ai/playbooks/`, says which one it chose in
one line, then runs it. Two exist:

- **`debug`** — reproduce, isolate, state the root cause in one sentence, make
  the smallest fix, keep the regression test, clean up failed experiments. It
  carries a production-incident branch that starts read-only and proposes
  mitigation (rollback, feature flag, scale up) before any proper fix.
- **`refactor`** — state what improves and which behaviour must not change, add
  characterization tests where coverage is missing, plan steps that each leave
  the build and tests green, checkpoint between steps, and revert a step that
  goes red when the cause is not clear within one attempt.

If no playbook fits, it says so and continues with the operating loop.

## Next

- [getting-started.md](getting-started.md) — install it.
- [extending.md](extending.md) — add your own.
