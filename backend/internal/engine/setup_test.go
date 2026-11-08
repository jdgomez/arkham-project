package engine_test

import (
	"testing"

	"github.com/jdgomez/arkham-project/backend/internal/engine"
)

func TestDemoGameStartsAtTurn1Upkeep(t *testing.T) {
	game := demoGame(t)

	if game.Turn != 1 {
		t.Errorf("turn = %d, want 1", game.Turn)
	}
	if game.Phase != engine.PhaseUpkeep {
		t.Errorf("phase = %q, want %q", game.Phase, engine.PhaseUpkeep)
	}
}

func TestDemoGameHasAFixedSeed(t *testing.T) {
	first := demoGame(t)
	second := demoGame(t)

	if first.Seed == 0 {
		t.Error("seed is zero, want a fixed non-zero seed")
	}
	if first.Seed != second.Seed {
		t.Errorf("seeds differ between demo games: %d and %d", first.Seed, second.Seed)
	}
}
