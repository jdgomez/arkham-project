# Architecture and decision record

Living record of the technical decisions that hold for the whole project, each with its reason and date.
Product scope, phases and glossary live in [`SESSION0.md`](../SESSION0.md). Milestone-specific decisions live in the milestone plans under [`docs/plans/`](plans/).

When a decision changes, update its entry (keep the old reason in one line if it helps) and change the date.

## Way of working

### D1 - Brain / Executor / Reviewer with branches and PRs
_2026-11-08_

Work is planned by Brain with the owner, implemented by an executor agent and validated by a separate reviewer agent. Each agent has its own goal, tools and responsibility, so no single agent can bend the checks to reach its own goal. Every change goes through a branch and a GitHub PR. The owner reads every PR and merges it with squash, so each PR is one easy rollback point on `main`.

No OpenSpec and no no-mistakes: this is a personal project without team or deploy-pipeline requirements.

### D2 - Planning artifacts
_2026-11-08_

- **Milestone plan** (`docs/plans/<milestone>.md`): goal, milestone-specific decisions and task list with acceptance criteria. The executor ticks tasks off as they are done.
- **Decision record** (this file): decisions that hold over the whole project.

Both are written by Brain with the owner and committed directly to `main` once the owner approves them in conversation (they have already been reviewed there, and the executor's branch is created from `main`).

### D3 - PR size, branches and commit messages
_2026-11-08_

- Several small PRs per milestone, each covering one coherent group of plan tasks and readable in one sitting.
- Squash commit titles follow [Conventional Commits](https://www.conventionalcommits.org/) (`feat:`, `fix:`, `test:`, `docs:`, `refactor:`, `chore:`), so the history of `main` reads as an index of features.
- Branch names match the type: `feat/<slug>`, `fix/<slug>`, `docs/<slug>`, ...

### D4 - Code readability
_2026-11-08_

The owner reads all code and is not yet an expert in Go. Code must be readable first.

- SOLID and Clean Code, applied the Go way: small interfaces defined where they are used, composition over hierarchies, single-responsibility packages, clear names, no speculative abstraction layers.
- Doc comments on every exported identifier (Go standard).
- Functions that implement a game rule cite the rule they apply, e.g. `// Rulebook: "Movement Phase" - ...`, so code can be checked against the rulebook.
- No line-by-line comments: names should make the code explain itself.

## Repository and tooling

### D5 - Single repository
_2026-11-08_

One repo with `backend/`, `frontend/` and `content/`. The Game State and Available Actions are a contract that changes on both sides at once: in one repo a feature is one PR and one squash commit, so a rollback reverts both sides together.

```
backend/    Go: engine, HTTP adapter, persistence
frontend/   Svelte 5 + TypeScript + Vite SPA
content/    YAML game content (cards, locations, monsters, ...)
docs/       architecture.md, plans/
Makefile    single entry point for humans, agents and CI
```

### D6 - Makefile as the single entry point
_2026-11-08_

The root `Makefile` holds every routine command, used the same way by the owner, the agents and CI. Expected targets: `make dev`, `make build`, `make test`, `make lint`, `make gen` (generate TS types). `make` was chosen over `just` because it needs no extra install.

### D7 - Lint and format
_2026-11-08_

- Back end: `gofmt` + `golangci-lint`, with linters that turn D4's mechanical rules into checks: doc comments on exported identifiers, magic numbers, cognitive complexity, function length and duplication (e.g. `revive`, `mnd`, `gocognit`, `funlen`, `dupl`).
- Front end: `eslint` + `prettier` + `svelte-check`, with the equivalent `eslint` rules (complexity, magic numbers).
- Configured in the golden slice. What no linter can express is checked by the reviewer through the project review skill (`.claude/skills/project-review/`), which changes in the same pull request as any decision that affects it.

### D8 - CI on GitHub Actions
_2026-11-08_

A simple workflow runs tests, lint and content validation on every PR. It is a check independent from the agents, and the owner sees the result on the PR before merging. There is no deploy pipeline.

### D9 - Versions and package manager
_2026-11-08_

Latest stable Go, Node LTS, npm (ships with Node; pnpm only pays off in repos with many front-end packages).

## Back end

### D10 - Engine as a pure library
_2026-11-08_

The rules engine (`backend/internal/engine`) knows nothing about HTTP, files or databases: it takes a Game State and an Action and returns the new Game State. The HTTP layer is a thin adapter on top. Rule tests stay simple and fast, rule code reads without noise, and changing the transport does not touch the rules.

### D11 - API: HTTP + JSON, action based
_2026-11-08_

- `GET /games/{id}`: Game State with its Available Actions.
- `POST /games/{id}/actions`: submit the chosen Action.

This is already RPC style (each Action is a command); rule complexity lives in the engine, not in the API shape. Revisit (gRPC / Connect) only when a concrete feature does not fit. WebSockets, if needed, come with remote multiplayer (Phase 4).

### D12 - Shared types generated from Go
_2026-11-08_

The Go structs are the source of truth for the contract. TypeScript types are generated from them (e.g. with `tygo`) via `make gen`, so the front end never uses a field the back end no longer sends.

### D13 - Content as YAML, loaded strictly
_2026-11-08_

Game content lives in `content/` as YAML, one file per content type, loaded by the back end at start-up. YAML reads better for long card texts in PR diffs. Strictness comes from the loader: unknown or misspelled fields are errors, and a validation test over all of `content/` (cross references exist, values in range) runs on every PR.

### D14 - Saved games as seed + action log
_2026-11-08_

A saved game is the seed plus the ordered list of Actions, together with the engine version. Loading replays the Actions. Files are small, any buggy game can be reproduced exactly (and turned into a scenario test), and it is the natural base for undo or replays. It also serves as the game log for debugging; add more logging only when a real need appears. Stored as files on disk, no database.

Risk: a rule change can make an old save incompatible; the stored engine version makes that detectable.

## Front end

### D15 - Svelte 5 SPA without SvelteKit
_2026-11-08_

Svelte 5 + TypeScript + Vite, no SvelteKit: there is already a Go back end and the game is a single stateful screen, so routing and server-side rendering are not needed. Fewer moving parts to understand in each PR.

### D16 - Dev and play modes
_2026-11-08_

- `make dev`: Go server + Vite dev server proxying `/api` to the back end.
- `make build`: a single Go binary that embeds the built front end and serves everything.

## Testing

### D17 - TDD and test levels
_2026-11-08_

- **TDD:** the executor works one acceptance criterion at a time: a test that fails for the right reason, the least code that makes it pass, then refactor while green - never all the tests of a slice first and the code after (`tdd` skill). The reviewer checks that the new tests fail without the change (red check) and that no test was weakened, skipped or deleted without justification.
- **Engine:** table-driven unit tests per rule, plus seed-based scenario tests that play a scripted sequence of Actions and check the resulting Game State. Scenario tests double as readable rule stories.
- **Content:** the validation test of D13.
- **Front end:** few component tests; Playwright E2E tests per milestone for the main player paths.
- **Manual testing:** anything not covered by functional or Playwright tests is tested by hand by the owner, who flags what is worth automating.
