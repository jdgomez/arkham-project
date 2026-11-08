# arkham-project

Product scope, phases and glossary: [`SESSION0.md`](SESSION0.md). Project-wide technical decisions: [`docs/architecture.md`](docs/architecture.md) (D1-D17). Milestone plans: `docs/plans/<milestone>.md`.

## Way of working

- **Methodology: TDD** (Test Driven Development). Agent personas: `tdd-executor` / `tdd-reviewer`. Plans are written with the `tdd-plan` skill; no OpenSpec, no no-mistakes (D1, D2).
- **Base branch:** `main`. Branches and squash titles follow D3 (Conventional Commits, `feat/<slug>`, `fix/<slug>`, ...). The executor opens the pull request; the owner merges with squash.
- **Routine commands** (D6, from the repo root): `make test`, `make lint`, `make gen`, `make build`, `make dev`.
- **Review rules** (tdd-reviewer): `.claude/skills/project-review/SKILL.md`, beyond the linters of D7. Update it in the same pull request as any decision that changes what should be reviewed.
- **Red check** (tdd-reviewer): `--cmd 'make test'` and `--tests '_test\.go$|/testdata/|\.test\.ts$|\.spec\.ts$'` until the golden slice settles the front-end test layout; update this line if it changes.
