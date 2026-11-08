package games

import (
	"fmt"

	"github.com/jdgomez/arkham-project/backend/internal/content"
	"github.com/jdgomez/arkham-project/backend/internal/engine"
)

// NewDemo returns a store holding only the demo game.
func NewDemo(loaded content.Content) (*Memory, error) {
	game, err := engine.DemoGame(loaded)
	if err != nil {
		return nil, fmt.Errorf("setting up the demo game: %w", err)
	}
	store := NewMemory()
	store.Put(engine.DemoGameID, game)
	return store, nil
}
