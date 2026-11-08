package engine

import (
	"fmt"

	"github.com/jdgomez/arkham-project/backend/internal/content"
)

// DemoGameID is the id of the demo game the server starts with.
const DemoGameID = "demo"

// demoSeed is the fixed seed of the demo game, so every run plays the same.
const demoSeed = 20261108

// demoAncientOne is the id of the Ancient One the demo game faces.
const demoAncientOne = "azathoth"

// NewGame sets up a game at the start of turn 1, with every investigator at
// their starting location, facing the Ancient One with the given id.
//
// Rulebook: "Turn Sequence" - a turn starts with the Upkeep phase.
func NewGame(loaded content.Content, seed int64, ancientOne string) (Game, error) {
	if _, ok := loaded.AncientOne(ancientOne); !ok {
		return Game{}, fmt.Errorf("setting up a game: Ancient One %q does not exist", ancientOne)
	}
	game := Game{Seed: seed, Turn: 1, Phase: PhaseUpkeep, AncientOne: ancientOne}
	for _, investigator := range loaded.Investigators {
		game.Investigators = append(game.Investigators,
			InvestigatorState{ID: investigator.ID, Location: investigator.StartingLocation})
	}
	return game, nil
}

// DemoGame is the fixture the server starts with and the tests build on.
func DemoGame(loaded content.Content) (Game, error) {
	return NewGame(loaded, demoSeed, demoAncientOne)
}
