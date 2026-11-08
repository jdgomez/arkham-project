package engine_test

import (
	"os"
	"testing"

	"github.com/jdgomez/arkham-project/backend/internal/content"
	"github.com/jdgomez/arkham-project/backend/internal/engine"
)

// loadContent loads the real content directory of the repository.
func loadContent(t *testing.T) content.Content {
	t.Helper()
	loaded, err := content.Load(os.DirFS("../../../content"))
	if err != nil {
		t.Fatalf("loading content: %v", err)
	}
	return loaded
}

func demoGame(t *testing.T) engine.Game {
	t.Helper()
	return engine.DemoGame(loadContent(t))
}
