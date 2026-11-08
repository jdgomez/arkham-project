package engine

import (
	"errors"
	"slices"
)

// ActionType says what an Action does.
type ActionType string

// ActionEndPhase ends the current phase. It is provisional: each phase gets
// its own Actions from milestone (b) on.
const ActionEndPhase ActionType = "End phase"

// Action is a choice the player makes.
type Action struct {
	Type ActionType `json:"type"`
}

// ErrActionNotAvailable is returned when an Action is not among the Available
// Actions of the game.
var ErrActionNotAvailable = errors.New("action not available")

// phaseOrder is the play order of a turn.
//
// Rulebook: "Turn Sequence".
var phaseOrder = []Phase{
	PhaseUpkeep,
	PhaseMovement,
	PhaseArkhamEncounters,
	PhaseOtherWorldEncounters,
	PhaseMythos,
}

// AvailableActions lists what the player can do in the game.
func AvailableActions(Game) []Action {
	return []Action{{Type: ActionEndPhase}}
}

// Apply returns the game after the Action. The given game is not changed.
func Apply(game Game, action Action) (Game, error) {
	if !slices.Contains(AvailableActions(game), action) {
		return game, ErrActionNotAvailable
	}
	return endPhase(game), nil
}

// endPhase moves to the next phase; after Mythos, to Upkeep of the next turn.
func endPhase(game Game) Game {
	next := slices.Index(phaseOrder, game.Phase) + 1
	if next == len(phaseOrder) {
		game.Turn++
		next = 0
	}
	game.Phase = phaseOrder[next]
	return game
}
