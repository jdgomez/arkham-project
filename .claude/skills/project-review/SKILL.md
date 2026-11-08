---
name: project-review
description: arkham-project review rules - what tdd-reviewer checks on every pull request of this project, beyond what its linters catch. Built from the decision record (docs/architecture.md, D1-D17).
---

# arkham-project review rules

Check the change against each rule. Report only what you can point at (file and line, which rule, why); zero findings is a valid result. A broken rule is blocking. Each rule cites the decision it comes from in `docs/architecture.md`.

## Architecture

- **The engine is a pure library (D10).** Nothing under `backend/internal/engine` imports HTTP, file, OS, database or YAML packages, or reads the clock or global randomness: it takes a Game State and an Action and returns the new Game State. Randomness comes from the game's seed.
- **The HTTP layer is a thin adapter (D10, D11).** Handlers decode the request, call the engine and encode the result. No rule logic, no rule-dependent branching in the transport layer.
- **The API stays action based (D11).** Game state is read with `GET /games/{id}` (state plus Available Actions) and changed only with `POST /games/{id}/actions`. A new endpoint needs a plan decision.
- **Shared types are generated (D12).** TypeScript types for the Game State, Actions and Available Actions come from the Go structs via `make gen`: no hand-written duplicates in the front end, and generated files are up to date with the Go structs in the same pull request.
- **Content is data, loaded strictly (D13).** Game content lives in `content/` as YAML, one file per content type; the loader rejects unknown or misspelled fields, and new content types are covered by the content validation test (cross references exist, values in range). No game content hardcoded in Go or TypeScript.
- **Saves are seed plus action log (D14).** A saved game stores the seed, the ordered Actions and the engine version - nothing derived from them.
- **The front end is a Svelte 5 SPA without SvelteKit (D15).** No SvelteKit, no server-side rendering, no client-side router added without a plan decision.

## Code

- **Rule functions cite the rulebook (D4).** A function that implements a game rule has a comment citing the rule it applies, e.g. `// Rulebook: "Movement Phase" - ...`.
- **Go the readable way (D4).** Small interfaces defined in the package that uses them, not next to the implementation; composition, not type hierarchies; one responsibility per package; no abstraction layer without a second real use or a test that needs it.
- **No line-by-line comments (D4).** Comments explain why or cite a rule; names explain what.

## Tests

- **Tests observe behavior at the seams (D17).** Tests only observe behavior at the seams listed in the milestone plan; a test reaching into internals (unexported helpers, private state, mocks of the project's own code) is a finding.
- **Engine rules (D17):** each rule has table-driven unit tests; behavior spanning several Actions is covered by a seed-based scenario test that plays a scripted sequence of Actions and checks the resulting Game State.
- **Front end (D17):** few component tests; a milestone's main player paths are covered by Playwright E2E tests before the milestone closes.
- **Content (D13):** the validation test runs over all of `content/`.

## Not rules here

- Formatting, import order, doc comments on exported identifiers, magic numbers, complexity and duplication: the linters enforce them (D7). Don't repeat them as findings.
- Naming taste, alternative designs or refactors no rule asks for: suggestions at most, never blocking.
