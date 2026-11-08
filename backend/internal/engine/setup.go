package engine

import "github.com/jdgomez/arkham-project/backend/internal/content"

// DemoGameID is the id of the demo game the server starts with.
const DemoGameID = "demo"

// demoSeed is the fixed seed of the demo game, so every run plays the same.
const demoSeed = 20261108

// NewGame sets up a game at the start of turn 1, with every investigator at
// their starting location and the first Ancient One of the content.
//
// Rulebook: "Turn Sequence" - a turn starts with the Upkeep phase.
func NewGame(loaded content.Content, seed int64) Game {
	game := Game{Seed: seed, Turn: 1, Phase: PhaseUpkeep, AncientOne: loaded.AncientOnes[0].ID}
	for _, investigator := range loaded.Investigators {
		game.Investigators = append(game.Investigators,
			InvestigatorState{ID: investigator.ID, Location: investigator.StartingLocation})
	}
	return game
}

// DemoGame is the fixture the server starts with and the tests build on.
func DemoGame(loaded content.Content) Game {
	return NewGame(loaded, demoSeed)
}
