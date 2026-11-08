# arkham-project

Product scope, phases and glossary: [`SESSION0.md`](SESSION0.md). Project-wide technical decisions: [`docs/architecture.md`](docs/architecture.md) (D1-D17). Milestone plans: `docs/plans/<milestone>.md`.

## Way of working

- **Methodology: TDD** (Test Driven Development). Agent personas: `tdd-executor` / `tdd-reviewer`. Plans are written with the `tdd-plan` skill; no OpenSpec, no no-mistakes (D1, D2).
- **Base branch:** `main`. Branches and squash titles follow D3 (Conventional Commits, `feat/<slug>`, `fix/<slug>`, ...). The executor opens the pull request; the owner merges with squash.
- **Routine commands** (D6, from the repo root): `make test`, `make lint`, `make gen`, `make build`, `make dev`.
- **Review rules** (tdd-reviewer): `.claude/skills/project-review/SKILL.md`, beyond the linters of D7. Update it in the same pull request as any decision that changes what should be reviewed.
- **`.claude/flow.json`** holds the machine-readable commands used by the agents' hooks and scripts: `base`, `test`, `lint` (run by the executor's done gate before a `DONE` is accepted) and `test_files` (the red check's pattern, provisional until the golden slice settles the front-end test layout). Keep it in step with the `Makefile`.
