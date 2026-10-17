# /arena - compare approaches before committing

Use when the choice is costly to reverse. Skip it when one option is obviously better.

## Steps

1. **Criteria.** Pick the criteria that matter for this task (for example latency, memory, simplicity, testability, operability, migration cost). Weight them if they conflict. Do not use a fixed list.
2. **Candidates.** Produce two approaches, three at most. They must differ in a real trade-off, not in cosmetics. Name each by what it is ("in-process cache", "Redis"), not "Alpha" and "Beta".
3. **Sketch.** For each: how it works in a few lines (not full code), main risks, and the cost of reversing it later.
4. **Compare.** One table, candidates as columns, criteria as rows.
5. **Recommend.** One pick, the reason, and what finding would change your mind.
6. **Stop.** Implement nothing until the user chooses.

If one candidate is clearly dominated, say so and drop it rather than padding the table.
