package engine_test

import (
	"os"
	"testing"

	"github.com/jdgomez/arkham-project/backend/internal/content"
	"github.com/jdgomez/arkham-project/backend/internal/contentyaml"
	"github.com/jdgomez/arkham-project/backend/internal/engine"
)

// loadContent loads the real content directory of the repository.
func loadContent(t *testing.T) content.Content {
	t.Helper()
	loaded, err := contentyaml.Load(os.DirFS("../../../content"))
	if err != nil {
		t.Fatalf("loading content: %v", err)
	}
	return loaded
}

func demoGame(t *testing.T) engine.Game {
	t.Helper()
	game, err := engine.DemoGame(loadContent(t))
	if err != nil {
		t.Fatalf("setting up the demo game: %v", err)
	}
	return game
}

func viewOf(t *testing.T, game engine.Game) engine.PlayerView {
	t.Helper()
	view, err := engine.ViewOf(loadContent(t), game)
	if err != nil {
		t.Fatalf("ViewOf returned error: %v", err)
	}
	return view
}
