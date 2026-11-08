// Package engine holds the rules of the game. It knows nothing about HTTP,
// files or databases: it takes a game state and returns a new one.
package engine

// Phase is one step of a turn.
type Phase string

// The phases of a turn, in play order.
//
// Rulebook: "Turn Sequence" - Upkeep, Movement, Arkham Encounters, Other World
// Encounters, Mythos.
const (
	PhaseUpkeep               Phase = "Upkeep"
	PhaseMovement             Phase = "Movement"
	PhaseArkhamEncounters     Phase = "Arkham Encounters"
	PhaseOtherWorldEncounters Phase = "Other World Encounters"
	PhaseMythos               Phase = "Mythos"
)

// Game is everything that changes during play. Content (the map, cards) is
// referenced by id and never copied here.
type Game struct {
	Seed          int64               `json:"seed"`
	Turn          int                 `json:"turn"`
	Phase         Phase               `json:"phase"`
	Investigators []InvestigatorState `json:"investigators"`
	AncientOne    string              `json:"ancientOne"`
}

// InvestigatorState is what changes about an investigator during play. The
// investigator's content is referenced by id.
type InvestigatorState struct {
	ID       string `json:"id"`
	Location string `json:"location"`
}
