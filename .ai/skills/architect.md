# /architect - design before code

Goal: agree on the design with the user before any implementation exists.

## Steps

1. **Frame.** Restate the goal and constraints in at most three lines. List assumptions. Ask only questions that block the design.
2. **Read.** Read the existing code the design touches, once. Note the conventions you must respect.
3. **Design.** Cover only the sections that apply:
   - Components and their responsibilities
   - Data model: entities, key fields, relations (Mermaid `erDiagram` when there are more than two entities)
   - Data flow (Mermaid `sequenceDiagram` or `flowchart`)
   - Interface contracts: signatures, inputs, outputs, error cases
   - Failure modes and what happens in each
   - Migration and rollout, if existing data or users are affected
   - Test strategy
4. **Decisions.** For each significant choice, name the alternative you rejected and why, in one line.
5. **Stop.** Write no implementation code. Ask for approval or changes.

## Output limits

One page. If it needs more, the scope is too big: propose a split.
