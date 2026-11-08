package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jdgomez/arkham-project/backend/internal/api"
	"github.com/jdgomez/arkham-project/backend/internal/engine"
	"github.com/jdgomez/arkham-project/backend/internal/games"
)

func newServer() http.Handler {
	store := games.NewMemory()
	store.Put("demo", engine.DemoGame())
	return api.NewHandler(store)
}

func get(t *testing.T, path string) *httptest.ResponseRecorder {
	t.Helper()
	response := httptest.NewRecorder()
	newServer().ServeHTTP(response, httptest.NewRequest(http.MethodGet, path, nil))
	return response
}

func TestGetDemoGameReturnsTheTurnAndPhase(t *testing.T) {
	response := get(t, "/games/demo")

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	var view struct {
		Turn  int    `json:"turn"`
		Phase string `json:"phase"`
	}
	if err := json.NewDecoder(response.Body).Decode(&view); err != nil {
		t.Fatalf("decoding body: %v", err)
	}
	if view.Turn != 1 || view.Phase != "Upkeep" {
		t.Errorf("view = turn %d, phase %q; want turn 1, phase \"Upkeep\"", view.Turn, view.Phase)
	}
}

func TestGetUnknownGameReturnsNotFound(t *testing.T) {
	response := get(t, "/games/unknown")

	if response.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", response.Code, http.StatusNotFound)
	}
}
