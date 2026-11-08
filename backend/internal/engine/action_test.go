package engine_test

import (
	"errors"
	"testing"

	"github.com/jdgomez/arkham-project/backend/internal/engine"
)

// Phase order: Rulebook, "Turn Sequence".
func TestEndPhaseMovesToTheNextPhase(t *testing.T) {
	tests := []struct {
		name  string
		start engine.Game
		want  engine.Game
	}{
		{
			name:  "Upkeep goes to Movement",
			start: engine.Game{Turn: 1, Phase: "Upkeep"},
			want:  engine.Game{Turn: 1, Phase: "Movement"},
		},
		{
			name:  "Movement goes to Arkham Encounters",
			start: engine.Game{Turn: 1, Phase: "Movement"},
			want:  engine.Game{Turn: 1, Phase: "Arkham Encounters"},
		},
		{
			name:  "Arkham Encounters goes to Other World Encounters",
			start: engine.Game{Turn: 1, Phase: "Arkham Encounters"},
			want:  engine.Game{Turn: 1, Phase: "Other World Encounters"},
		},
		{
			name:  "Other World Encounters goes to Mythos",
			start: engine.Game{Turn: 1, Phase: "Other World Encounters"},
			want:  engine.Game{Turn: 1, Phase: "Mythos"},
		},
		{
			name:  "Mythos goes to Upkeep of the next turn",
			start: engine.Game{Turn: 1, Phase: "Mythos"},
			want:  engine.Game{Turn: 2, Phase: "Upkeep"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := engine.Apply(test.start, engine.Action{Type: "End phase"})

			if err != nil {
				t.Fatalf("Apply returned error: %v", err)
			}
			if got != test.want {
				t.Errorf("game = %+v, want %+v", got, test.want)
			}
		})
	}
}

func TestApplyRejectsAnActionThatIsNotAvailable(t *testing.T) {
	start := engine.DemoGame()

	got, err := engine.Apply(start, engine.Action{Type: "Dance"})

	if !errors.Is(err, engine.ErrActionNotAvailable) {
		t.Errorf("error = %v, want ErrActionNotAvailable", err)
	}
	if got != start {
		t.Errorf("game = %+v, want it unchanged: %+v", got, start)
	}
}

func TestTheOnlyAvailableActionOfTheDemoGameIsEndPhase(t *testing.T) {
	view := engine.ViewOf(engine.DemoGame())

	want := []engine.Action{{Type: "End phase"}}
	if len(view.AvailableActions) != 1 || view.AvailableActions[0] != want[0] {
		t.Errorf("available actions = %+v, want %+v", view.AvailableActions, want)
	}
}
