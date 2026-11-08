# Milestone A: front/back skeleton and Game State

- **Status:** approved
- **Mode:** golden slice
- **Goal:** the player opens the app, sees the demo game (the map of the 2 districts, Mandy Thompson at the Library, Azathoth, turn and phase) and advances it with "End phase" through every phase of a turn, turn after turn.

## Decisions for this milestone

- **Demo fixture instead of "New game".** The server starts with one in-memory game with id `demo`, built by the engine's own setup function from the loaded content and a fixed seed - never a hand-written JSON. Tests (engine scenarios and E2E) start from the same fixture. It lives in memory and is reset when the server restarts. The front end shows the `demo` game at `/`.
- **Game state vs. content.** The game state holds only what changes during play (seed, turn, phase, each investigator's location, the Ancient One) and references content by id; it never copies it. Content (the map) is static, loaded once from `content/` (D13), the same for the demo and for a real game.
- **Player view.** `GET /games/{id}` returns the player's view: a pure engine function composes content + game state + Available Actions. The HTTP handler only encodes it (D10, D11). This is also where the per-player projection of Phase 4 will live.
- **Fixture contents:** turn 1, phase Upkeep, a fixed seed (stored in the state from the start, D14); Mandy Thompson at the Library, with no stats, cards or abilities yet; Azathoth as Ancient One, name only (no doom track).
- **Map:** Merchant District (Unvisited Isle, River Docks, The Unnamable) and Miskatonic U. (Administration, Library, Science Building), each location connected to its district's street, and the two streets connected to each other - a real connection on the board. The MVP adaptation of SESSION0 concerns only the monster movement arrows and is settled in milestone (d).
- **"End phase" is provisional:** the only Available Action in this milestone. It is replaced phase by phase as each phase gets its own rules from milestone (b) on.
- **Phase order** (Rulebook, "Turn Sequence"): Upkeep -> Movement -> Arkham Encounters -> Other World Encounters -> Mythos, then Upkeep of the next turn.

## Seams

| Seam | Tests here catch | Tests here miss |
|---|---|---|
| Engine: apply an Action to a game state; compose the player's view (state + Available Actions) | Rules: phase order, turn change, illegal Actions rejected; the content of the view | JSON, HTTP, UI |
| Content loading: load the `content/` directory | Unknown or misspelled fields, broken cross references, values out of range | Whether the content is faithful to the game (the owner checks) |
| HTTP API: `GET /games/{id}`, `POST /games/{id}/actions`, with the real engine, no mocks | JSON shape, status codes, illegal Action -> 4xx, unknown game -> 404 | Rule details (engine), UI |
| Svelte components: few, only where there is presentation logic | That Available Actions are rendered and clicking one emits it | Integration with the back end |
| E2E Playwright against the `make build` binary | The full trace: generated types, embedded front end, the player sees the state and acts on it | Rule details; slow, so few and only the main player paths |

No test is written outside these seams. `make gen` leaving no diff is a CI check, not a seam.

**Reference style per level** (the executor writes the first example of each in S1-S3, and later tests copy it):
- Engine rules: table-driven unit tests, one table per rule, each row a named case with input state, Action and expected state.
- Engine scenarios: start from the demo fixture, play a scripted list of Actions, check the resulting state.
- Content: one validation test over all of `content/`, plus small invalid inputs for the loader's error cases.
- HTTP: request in, status code and decoded JSON out, through the real router and engine.
- Component: render with given props, assert on what the player sees, interact as the player would.
- E2E: one test per main player path, written as the player's steps and what they see.

## Slices

### S1: Golden trace - the player sees the turn and the phase

- **Story:** the player opens the app -> sees "Turn 1 - Upkeep" of the demo game.
- **Class:** behavior (includes the project scaffolding: Makefile, Go module, Svelte/Vite app, linters, CI, type generation, embedded build, Playwright)
- **Blocked by:** none
- **Acceptance criteria:**
  - [x] Given the demo setup, then the game state is turn 1, phase Upkeep, with a fixed seed (seam: engine)
  - [x] Given the demo game, when `GET /games/demo`, then 200 with turn 1 and phase "Upkeep" (seam: HTTP API)
  - [x] When `GET /games/unknown`, then 404 (seam: HTTP API)
  - [x] Given the built binary, when the player opens `/`, then they see "Turn 1 - Upkeep" (seam: E2E)
  - [x] `make test`, `make lint`, `make gen`, `make build` and `make dev` work from the repo root (D6, D16)
  - [x] The D7 linters are configured for back end and front end
  - [x] CI runs tests, lint and a check that `make gen` leaves no diff, on every pull request (D8)
  - [x] `.claude/flow.json` matches the Makefile, and its `test_files` matches the real test layout (front-end tests included)

### S2: End phase - the player advances phases and turns

- **Story:** the player clicks "End phase" -> the game moves to the next phase, and after Mythos to Upkeep of the next turn.
- **Class:** behavior
- **Blocked by:** S1
- **Acceptance criteria:**
  - [x] Given each phase, when "End phase" is applied, then: Upkeep -> Movement, Movement -> Arkham Encounters, Arkham Encounters -> Other World Encounters, Other World Encounters -> Mythos, Mythos -> Upkeep with turn + 1 (seam: engine)
  - [x] Given the demo game, then the only Available Action is "End phase" (seam: engine)
  - [x] Given the demo game, when "End phase" is applied 5 times, then turn 2, Upkeep; 10 times, then turn 3, Upkeep (seam: engine, scenario)
  - [x] When `POST /games/demo/actions` with "End phase", then 200 with turn 1, phase Movement (seam: HTTP API)
  - [x] When `POST` an Action that is not in the Available Actions, then 4xx and the game state is unchanged (seam: HTTP API)
  - [x] Given a list of Available Actions, the component renders one button per Action, and clicking one emits that Action (seam: component)
  - [x] Given the app, when the player clicks "End phase", then they see "Turn 1 - Movement"; after 4 more clicks, "Turn 2 - Upkeep" (seam: E2E)

### S3: The map from content

- **Story:** the player opens the app -> sees the map of the 2 districts with their connections, Mandy Thompson at the Library and Azathoth as Ancient One.
- **Class:** behavior
- **Blocked by:** S2
- **Acceptance criteria:**
  - [ ] Given `content/`, when loaded, then 6 locations and 2 streets: each location connected to its district's street, and Merchant District Streets connected to Miskatonic U. Streets (seam: content loading)
  - [ ] Given a content file with an unknown field, when loaded, then the load fails naming the file and the field (seam: content loading)
  - [ ] Given a connection to an area that does not exist, then validation fails naming it (seam: content loading)
  - [ ] Given the demo game, then the player's view includes the board (8 areas and their connections), Mandy Thompson at the Library, and Azathoth as Ancient One (seam: engine)
  - [ ] Given the app, when the player opens `/`, then they see the 8 areas with their connections, Mandy Thompson at the Library and Azathoth (seam: E2E)

## Out of this milestone

- "New game" and `POST /games` (a D11 decision when it is needed).
- Movement, encounters, cards, investigator stats and abilities, the second investigator (milestone b).
- Doom track, gates, monsters and their movement arrows (milestones c, d).
- Saving games to disk (milestone e, D14).

## Open questions

- None blocking.
