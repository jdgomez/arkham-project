// Package engine holds the rules of the game. It knows nothing about HTTP,
// files or databases: it takes a game state and returns a new one.
package engine

// Phase is one step of a turn.
type Phase string

// Upkeep is the first phase of every turn.
const Upkeep Phase = "Upkeep"

// Game is everything that changes during play. Content (the map, cards) is
// referenced by id and never copied here.
type Game struct {
	Seed  int64 `json:"seed"`
	Turn  int   `json:"turn"`
	Phase Phase `json:"phase"`
}
