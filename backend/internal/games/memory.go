// Package games keeps the games that are being played.
package games

import (
	"sync"

	"github.com/jdgomez/arkham-project/backend/internal/engine"
)

// Memory stores games in memory: they are lost when the server restarts.
type Memory struct {
	mu    sync.RWMutex
	games map[string]engine.Game
}

// NewMemory returns an empty store.
func NewMemory() *Memory {
	return &Memory{games: map[string]engine.Game{}}
}

// Put stores the game under the id, replacing any previous one.
func (m *Memory) Put(id string, game engine.Game) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.games[id] = game
}

// Get returns the game with the id, and whether it exists.
func (m *Memory) Get(id string) (engine.Game, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	game, ok := m.games[id]
	return game, ok
}
