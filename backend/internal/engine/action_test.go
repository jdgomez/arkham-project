package engine_test

import (
	"errors"
	"reflect"
	"slices"
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
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("game = %+v, want %+v", got, test.want)
			}
		})
	}
}

func TestApplyRejectsAnActionThatIsNotAvailable(t *testing.T) {
	tests := []struct {
		name   string
		action engine.Action
	}{
		{"an unknown kind of Action", engine.Action{Type: "Dance"}},
		{"an Action without a type", engine.Action{}},
		{"the right type in the wrong case", engine.Action{Type: "end phase"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			start := demoGame(t)

			got, err := engine.Apply(start, test.action)

			if !errors.Is(err, engine.ErrActionNotAvailable) {
				t.Errorf("error = %v, want ErrActionNotAvailable", err)
			}
			if !reflect.DeepEqual(got, start) {
				t.Errorf("game = %+v, want it unchanged: %+v", got, start)
			}
		})
	}
}

func TestApplyFailsOnAPhaseThatDoesNotExist(t *testing.T) {
	start := engine.Game{Turn: 1, Phase: "Teatime"}

	got, err := engine.Apply(start, engine.Action{Type: "End phase"})

	if err == nil || errors.Is(err, engine.ErrActionNotAvailable) {
		t.Errorf("error = %v, want an error about the unknown phase", err)
	}
	if !reflect.DeepEqual(got, start) {
		t.Errorf("game = %+v, want it unchanged: %+v", got, start)
	}
}

func TestAvailableActions(t *testing.T) {
	tests := []struct {
		name string
		game engine.Game
		want []engine.Action
	}{
		{"the demo game", demoGame(t), []engine.Action{{Type: "End phase"}}},
		{"Mythos of a later turn", engine.Game{Turn: 7, Phase: "Mythos"}, []engine.Action{{Type: "End phase"}}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := engine.AvailableActions(test.game)

			if !slices.Equal(got, test.want) {
				t.Errorf("available actions = %+v, want %+v", got, test.want)
			}
		})
	}
}
