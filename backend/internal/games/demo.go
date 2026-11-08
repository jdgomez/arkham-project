package games

import (
	"github.com/jdgomez/arkham-project/backend/internal/content"
	"github.com/jdgomez/arkham-project/backend/internal/engine"
)

// NewDemo returns a store holding only the demo game.
func NewDemo(loaded content.Content) *Memory {
	store := NewMemory()
	store.Put(engine.DemoGameID, engine.DemoGame(loaded))
	return store
}
