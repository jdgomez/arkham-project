package engine_test

import (
	"testing"

	"github.com/jdgomez/arkham-project/backend/internal/content"
	"github.com/jdgomez/arkham-project/backend/internal/engine"
)

// smallContent is content built in Go: one location, one investigator and the
// given Ancient Ones.
func smallContent(t *testing.T, ancientOnes ...content.AncientOne) content.Content {
	t.Helper()
	loaded, err := content.New(
		[]content.Area{{ID: "library", Name: "Library", Connections: []string{"streets"}}},
		[]content.Area{{ID: "streets", Name: "Streets", Connections: []string{"library"}}},
		[]content.Investigator{{ID: "mandy", Name: "Mandy", StartingLocation: "library"}},
		ancientOnes,
	)
	if err != nil {
		t.Fatalf("building content: %v", err)
	}
	return loaded
}

func TestNewGameFacesTheAncientOneWithTheGivenId(t *testing.T) {
	loaded := smallContent(t,
		content.AncientOne{ID: "azathoth", Name: "Azathoth"},
		content.AncientOne{ID: "cthulhu", Name: "Cthulhu"},
	)

	game, err := engine.NewGame(loaded, 1, "cthulhu")

	if err != nil {
		t.Fatalf("NewGame returned error: %v", err)
	}
	if game.AncientOne != "cthulhu" {
		t.Errorf("ancient one = %q, want %q", game.AncientOne, "cthulhu")
	}
}

func TestNewGameFailsOnAnAncientOneThatDoesNotExist(t *testing.T) {
	loaded := smallContent(t, content.AncientOne{ID: "azathoth", Name: "Azathoth"})

	_, err := engine.NewGame(loaded, 1, "cthulhu")

	if err == nil {
		t.Error("NewGame returned no error, want one about the unknown Ancient One")
	}
}

func TestViewOfFailsOnContentThatTheGameReferencesButDoesNotExist(t *testing.T) {
	loaded := smallContent(t, content.AncientOne{ID: "azathoth", Name: "Azathoth"})
	tests := []struct {
		name string
		game engine.Game
	}{
		{"an investigator", engine.Game{
			Investigators: []engine.InvestigatorState{{ID: "ghost", Location: "library"}},
			AncientOne:    "azathoth",
		}},
		{"the Ancient One", engine.Game{AncientOne: "cthulhu"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := engine.ViewOf(loaded, test.game)

			if err == nil {
				t.Error("ViewOf returned no error, want one about the missing content")
			}
		})
	}
}
