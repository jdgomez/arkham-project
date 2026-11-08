package engine_test

import (
	"testing"

	"github.com/jdgomez/arkham-project/backend/internal/engine"
)

// play applies the Actions in order, starting from the demo game.
func play(t *testing.T, actions ...engine.Action) engine.Game {
	t.Helper()
	game := engine.DemoGame()
	for i, action := range actions {
		var err error
		game, err = engine.Apply(game, action)
		if err != nil {
			t.Fatalf("action %d (%q): %v", i+1, action.Type, err)
		}
	}
	return game
}

func endPhases(count int) []engine.Action {
	actions := make([]engine.Action, count)
	for i := range actions {
		actions[i] = engine.Action{Type: "End phase"}
	}
	return actions
}

func TestEndingPhasesAdvancesTurnAfterTurn(t *testing.T) {
	tests := []struct {
		name      string
		endPhases int
		wantTurn  int
		wantPhase engine.Phase
	}{
		{"a whole turn", 5, 2, "Upkeep"},
		{"two whole turns", 10, 3, "Upkeep"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			game := play(t, endPhases(test.endPhases)...)

			if game.Turn != test.wantTurn || game.Phase != test.wantPhase {
				t.Errorf("game = turn %d, %q; want turn %d, %q",
					game.Turn, game.Phase, test.wantTurn, test.wantPhase)
			}
		})
	}
}
