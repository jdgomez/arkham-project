package engine_test

import (
	"testing"

	"github.com/jdgomez/arkham-project/backend/internal/engine"
)

func TestDemoGameStartsAtTurn1Upkeep(t *testing.T) {
	game := engine.DemoGame()

	if game.Turn != 1 {
		t.Errorf("turn = %d, want 1", game.Turn)
	}
	if game.Phase != engine.Upkeep {
		t.Errorf("phase = %q, want %q", game.Phase, engine.Upkeep)
	}
}

func TestDemoGameHasAFixedSeed(t *testing.T) {
	first := engine.DemoGame()
	second := engine.DemoGame()

	if first.Seed == 0 {
		t.Error("seed is zero, want a fixed non-zero seed")
	}
	if first.Seed != second.Seed {
		t.Errorf("seeds differ between demo games: %d and %d", first.Seed, second.Seed)
	}
}
