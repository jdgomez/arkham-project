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
	return api.NewHandler(games.NewDemo())
}

func get(t *testing.T, path string) *httptest.ResponseRecorder {
	t.Helper()
	response := httptest.NewRecorder()
	newServer().ServeHTTP(response, httptest.NewRequest(http.MethodGet, path, nil))
	return response
}

func TestGetDemoGameReturnsTheTurnAndPhase(t *testing.T) {
	response := get(t, "/games/"+engine.DemoGameID)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if contentType := response.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Content-Type = %q, want %q", contentType, "application/json")
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
