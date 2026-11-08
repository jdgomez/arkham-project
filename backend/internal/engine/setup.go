package engine

// DemoGameID is the id of the demo game the server starts with.
const DemoGameID = "demo"

// demoSeed is the fixed seed of the demo game, so every run plays the same.
const demoSeed = 20261108

// NewGame sets up a game at the start of turn 1.
//
// Rulebook: "Turn Sequence" - a turn starts with the Upkeep phase.
func NewGame(seed int64) Game {
	return Game{Seed: seed, Turn: 1, Phase: PhaseUpkeep}
}

// DemoGame is the fixture the server starts with and the tests build on.
func DemoGame() Game {
	return NewGame(demoSeed)
}
