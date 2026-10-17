# /tdd - red, green, refactor

## Cycle (one behavior per cycle)

1. **Red.** Write one test for the next behavior. Run it. Confirm it fails for the right reason (missing behavior, not a typo or import error).
2. **Green.** Write the least code that makes it pass. Run the test, then the surrounding suite.
3. **Refactor.** Improve names and structure with tests green. Run the suite again.
4. Repeat for the next behavior.

## Rules

- No production code without a failing test that demands it.
- Bug fixes start with a regression test that reproduces the bug.
- Test behavior through public interfaces, not private details.
- Do not mock what you own unless there is no cheaper way. Mock at system boundaries (network, clock, filesystem).
- Never edit a test to make it pass unless the requirement itself changed. If it changed, say so.
- If the repo has no test harness, say so and propose the smallest one that fits the stack. Do not silently skip tests.
- Follow the repo's existing test framework, layout and naming.
