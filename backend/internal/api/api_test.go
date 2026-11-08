package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jdgomez/arkham-project/backend/internal/api"
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
	response := get(t, "/games/demo")

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

type view struct {
	Turn             int    `json:"turn"`
	Phase            string `json:"phase"`
	AvailableActions []struct {
		Type string `json:"type"`
	} `json:"availableActions"`
}

func decodeView(t *testing.T, response *httptest.ResponseRecorder) view {
	t.Helper()
	var decoded view
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil {
		t.Fatalf("decoding body: %v", err)
	}
	return decoded
}

func TestGetDemoGameListsEndPhaseAsTheOnlyAvailableAction(t *testing.T) {
	decoded := decodeView(t, get(t, "/games/demo"))

	if len(decoded.AvailableActions) != 1 || decoded.AvailableActions[0].Type != "End phase" {
		t.Errorf("available actions = %+v, want only \"End phase\"", decoded.AvailableActions)
	}
}

func post(handler http.Handler, path, body string) *httptest.ResponseRecorder {
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodPost, path, strings.NewReader(body)))
	return response
}

func TestPostEndPhaseReturnsTheNewState(t *testing.T) {
	response := post(newServer(), "/games/demo/actions", `{"type":"End phase"}`)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	decoded := decodeView(t, response)
	if decoded.Turn != 1 || decoded.Phase != "Movement" {
		t.Errorf("view = turn %d, phase %q; want turn 1, phase \"Movement\"", decoded.Turn, decoded.Phase)
	}
}

func TestPostAnActionThatIsNotAvailableIsRejectedAndTheGameIsUnchanged(t *testing.T) {
	server := newServer()

	response := post(server, "/games/demo/actions", `{"type":"Dance"}`)

	if response.Code < 400 || response.Code >= 500 {
		t.Errorf("status = %d, want 4xx", response.Code)
	}
	after := httptest.NewRecorder()
	server.ServeHTTP(after, httptest.NewRequest(http.MethodGet, "/games/demo", nil))
	decoded := decodeView(t, after)
	if decoded.Turn != 1 || decoded.Phase != "Upkeep" {
		t.Errorf("game = turn %d, phase %q; want it unchanged: turn 1, phase \"Upkeep\"", decoded.Turn, decoded.Phase)
	}
}

func TestPostToAnUnknownGameReturnsNotFound(t *testing.T) {
	response := post(newServer(), "/games/unknown/actions", `{"type":"End phase"}`)

	if response.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", response.Code, http.StatusNotFound)
	}
}

func TestGetUnknownGameReturnsNotFound(t *testing.T) {
	response := get(t, "/games/unknown")

	if response.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", response.Code, http.StatusNotFound)
	}
}
