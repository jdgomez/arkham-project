// Package games keeps the games that are being played.
package games

import (
	"errors"
	"sync"

	"github.com/jdgomez/arkham-project/backend/internal/engine"
)

// ErrNotFound is returned when there is no game with the id.
var ErrNotFound = errors.New("game not found")

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

// Update replaces the game with the id by what change returns, as one atomic
// step: no other Get or Update sees or interleaves with it. If change fails,
// the game stays as it was and the error is returned.
func (m *Memory) Update(id string, change func(engine.Game) (engine.Game, error)) (engine.Game, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	game, ok := m.games[id]
	if !ok {
		return engine.Game{}, ErrNotFound
	}
	next, err := change(game)
	if err != nil {
		return game, err
	}
	m.games[id] = next
	return next, nil
}
