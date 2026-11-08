# SESSION0 - Arkham Project
_Date: 2026-11-08_

## Problem
Arkham Horror 2nd Edition (Fantasy Flight Games) is out of print. The owner's physical copy (base game plus almost every expansion) is effectively irreplaceable, and the game is slow to set up and tear down, so it no longer gets played. The project is a faithful digital version of the game so it can be played comfortably from time to time, and shown to the owner's nephews in the future, when finding a physical copy will likely be impossible.

What "comfortable" means here: no setup or teardown, and every piece of bookkeeping that can be automated is automated. The fun is in the decisions (how to spend clues, when to risk sanity, where to go), not in tracking numbers by hand.

## Target user
- **Now:** the owner alone, playing single player and controlling several investigators. He knows the game and its rules well, but has not played it in 2-3 years, so the engine must enforce the rules without needing the physical rulebook.
- **Later:** a few friends in another city, playing remotely through a very simple login.
- Strictly private, personal use. The game is not licensed, so this is never published or distributed publicly.

## Vision
A web game with a "video game" feel (menus and panels, not a virtual tabletop) that is as faithful to the original rules as possible. The engine applies every rule it can: phases, dice rolls, card effects, encounters, victory and defeat. The player only makes decisions.

The front end and back end are strictly separated. The back end returns the Game State as JSON, including the Available Actions for the connected player. The front end only renders that state and sends back the chosen action. It never decides what is legal.

Success, in the long run, means being able to play a full game of AH2e with friends without opening the rulebook.

## Scope

### In
- Single player: one player controls several investigators (1-2 in Phase 1, up to 8 in Phase 3).
- A rules engine that automates every rule and every piece of bookkeeping it can.
- Content (cards, locations, investigators, Ancient Ones, monsters) defined as data files, so expansion sets can be added without rewriting the engine.
- A web UI, with the front end and back end separated by a JSON Game State API.
- Saved games (games are very long).
- A reproducible random seed for dice and shuffles (to replay exact games when debugging).
- Card text transcribed from the Arkham Horror Wiki. Anything missing is typed in by the owner from his physical copy.
- A schematic UI at first (cards as text, simple icons, the map as a graph of streets and locations).
- The full base game (Phase 3).

### Out
- Public release or distribution of any kind.
- Licensing or commercial use.

### Pending
- Remote multiplayer with friends and a minimal login system (Phase 4).
- Expansion sets (Phase 5).
- Original or AI-generated artwork as a swappable "skin" over the schematic UI. The source is undecided: there is no scanner yet, and a small paid AI image generation license is an option.
- Private deployment (likely the owner's Raspberry Pi behind a Cloudflare Tunnel/Access) once multiplayer arrives.

## Glossary
| Term | Agreed definition |
|------|-------------------|
| Engine | The program that knows and applies the AH2e rules: phases, dice rolls, card effects, encounters, victory and defeat. It does every piece of bookkeeping that can be automated, so the player only makes decisions. Separate from the UI. |
| Expansion set | A content package (cards, boards, a few new mechanics) added on top of the base game. Design requirement: content is data, so expansion sets can be added without rewriting the engine. |
| Player | A real person playing the game. In the future, a user with a login. One player can control several investigators. |
| Investigator | A game character (e.g. Mandy Thompson) controlled by a player. |
| Ancient One | The Great Old One the game is played against. It sets the doom track, the final battle and its special effects. Phase 1 uses Azathoth. |
| Game State | The JSON document the back end returns with everything a given player's screen needs: board, investigators, Ancient One, current phase, etc. |
| Available Actions | The part of the Game State that lists the decisions that player can legally make right now. The front end only shows them and sends back the choice. |
| Vertical slice | A version of the game that goes through every layer and phase (front end, back end, engine, Game State, full phase cycle) with a minimal subset of content and only the simplest variant of each mechanic. It validates the architecture before growing in breadth. |
| MVP adaptation | A rule from the original that is deliberately changed in the MVP because the content is reduced (e.g. the street connections for monster movement). Each one is documented and reverted when the content grows. |
| Action | A concrete order the player picks from the Available Actions; the back end validates it and applies it to the Game State. |
| Milestone plan | `docs/plans/<milestone>.md`: the milestone goal, milestone-specific decisions and the task list with acceptance criteria, ticked off by the executor. |
| Decision record | `docs/architecture.md`: the technical decisions that hold for the whole project, each with its reason and date. |

## Phases

### Phase 1 - Vertical slice (MVP)
The smallest complete game: every phase and layer, with minimal content and only the simplest mechanics.

- **Stack:** Go back end, Svelte + TypeScript front end (Vite). Game State and Available Actions as JSON. Content as data files.
- **Players:** single player, 1-2 investigators: Mandy Thompson (an active ability that forces a decision mid-roll: re-roll failed dice) and Michael McGlen (a passive ability: reduce Stamina loss). Exact texts to be checked on the wiki.
- **Ancient One:** Azathoth.
- **Map:** 2 districts only: Merchant District (Unvisited Isle, River Docks, The Unnamable) and Miskatonic U. (Administration, Library, Science Building), each with its street. MVP adaptation: the streets are connected so monsters can keep moving between them instead of getting stuck (monster movement arrows do not form a loop between these two districts). To be detailed when monster movement is designed.
- **Arkham encounters:** 5 per location. An encounter is just "draw a card and resolve it". No location special abilities.
- **Other Worlds:** 2 (working proposal: City of the Great Race and R'lyeh), 5 encounters each.
- **Mythos:** 10 Headline cards only, whose gate locations fall on the unstable locations of the 2 districts (real cards first, reassigned only if there are not enough).
- **Monsters:** 3, the simplest only (no fast, flying or stationary monsters).
- **Investigator cards:** 10 Common Items, 10 Spells, 10 Allies, 10 Skills.
- **Excluded from Phase 1:** Environments, Rumors, Unique Items, Blessing/Curse, Retainer, Bank Loan, Deputy, Silver Twilight Lodge, Terror Track, Outskirts, Monster Limit, Sky, Lost in Time and Space, location special abilities.
- **Persistence:** saved games and a reproducible random seed.
- **UI:** schematic.

**Success criterion:** the engine plays games that go through every phase (Upkeep, Movement, Arkham Encounters, Other World Encounters, Mythos) turn after turn, reshuffling decks when they run out, until a legal ending (victory by closing or sealing gates, or defeat by Azathoth awakening), without opening the rulebook.

Working breakdown into milestones: (a) front/back skeleton and Game State; (b) investigators, movement and Arkham encounters; (c) gates, Other Worlds and Mythos; (d) monsters and combat; (e) victory, defeat and saved games.

### Phase 2 - Special rules sampler
Still on the vertical slice, and still with 1-2 investigators: add one example of each special rule still missing, one at a time, so each detail can be tracked and verified on its own. For example:
- Monster Limit (and the Outskirts / Terror Track mechanics it brings with it).
- One fast monster, one flying monster (and the Sky), one stationary monster.
- Environment Mythos cards, and a single Rumor.
- Two or three Unique Items.
- Locations with special abilities, one by one (e.g. a store, the Science Building's ability).
- The rest of the special mechanics (Blessing/Curse, Retainer, Bank Loan, Deputy, Silver Twilight Lodge, Lost in Time and Space...).

Success criterion: every rule mechanic of the base game exists and works with at least one piece of content.

### Phase 3 - Full base game
Add the remaining content on top of the existing rules: all 9 districts and the real map (MVP adaptations reverted), all encounters and Mythos cards, every monster and investigator card, all 16 investigators, all 8 Ancient Ones with their final battles, up to 8 investigators. New rule exceptions should be minimal (maybe a few investigator abilities).

### Phase 4 - Remote multiplayer
A minimal login, several players in one game (each controlling their own investigators), push updates through WebSockets, private deployment. The Game State is designed "per player" from Phase 1 so this does not require a rewrite.

### Phase 5 - Expansion sets
Add expansion sets one at a time, starting with the one the owner picks.

## Constraints and technical decisions
The full, living decision record (way of working, repo layout, tooling, API, persistence, testing) is [`docs/architecture.md`](docs/architecture.md). The items below are the product-level constraints from the first Session 0 pass.

- **Back end: Go.** Static typing suits a large rules state machine. It compiles to a single dependency-free binary (also for ARM, so it can be deployed on a Raspberry Pi), and has mature HTTP and WebSocket support. TypeScript types for the Game State are generated from the Go structs, so the front end shares the contract.
- **Front end: TypeScript + Svelte (Vite).** The front end is "render a Game State". A reactive framework handles DOM updates, which would be fragile to hand-write in vanilla TS as panels, modals and dice animations grow. Svelte is the closest to vanilla (plain HTML/CSS/TS components, compiled, small runtime). The owner prefers TypeScript and knows little about frameworks. Svelte was the recommendation and it was accepted.
- **Strict front/back separation:** the back end owns all rules and legality. The front end renders the Game State and submits one of the Available Actions.
- **Content as data** (YAML files loaded strictly by the engine, see D13), as a requirement for expansion sets.
- **Saved games** from Phase 1, with a reproducible seed for all randomness (stored as seed + action log, see D14).
- **Language:** everything official and everything in the repo (code, docs, commits, domain terms) is in English.
- **Card text source:** Arkham Horror Wiki, plus manual transcription from the physical copy.
- **No licence:** private, personal use only.
- **Repo:** `~/Desktop/Projects/arkham-project`, remote `git@github.com:jdgomez/arkham-project.git`.

## Open questions
- Exact choice of the 2 Other Worlds (working proposal: City of the Great Race and R'lyeh) and how Other World encounters, which are drawn by gate colour from a shared deck, map onto "5 encounters per world".
- The exact MVP adaptation for monster movement between the 2 districts. To be discussed in detail with the owner when monster movement is designed.
- Which of the 6 locations are unstable (gate sites). To be checked on the wiki, together with which real Headline cards open gates there.
- Location special abilities: Phase 1 originally kept them for its end, and the final phasing moves them to Phase 2 (one by one). Assumed Phase 2. To confirm.
- Some Phase 2 examples (e.g. a store) live in districts outside the 2 MVP districts. Phase 2 may need to add the odd location or district.
- How Monster Limit, Outskirts and Terror Track are introduced in Phase 2 (together, or one at a time).
- Artwork source (scans, AI generation) and when to introduce it.
- Exact texts of Mandy Thompson's and Michael McGlen's abilities, to be verified on the wiki.
