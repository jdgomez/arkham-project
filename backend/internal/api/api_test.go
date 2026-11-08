package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/jdgomez/arkham-project/backend/internal/api"
	"github.com/jdgomez/arkham-project/backend/internal/contentyaml"
	"github.com/jdgomez/arkham-project/backend/internal/engine"
	"github.com/jdgomez/arkham-project/backend/internal/games"
)

func newServer(t *testing.T) http.Handler {
	t.Helper()
	loaded, err := contentyaml.Load(os.DirFS("../../../content"))
	if err != nil {
		t.Fatalf("loading content: %v", err)
	}
	store, err := games.NewDemo(loaded)
	if err != nil {
		t.Fatalf("setting up the demo game: %v", err)
	}
	return api.NewHandler(store, loaded)
}

func get(t *testing.T, path string) *httptest.ResponseRecorder {
	t.Helper()
	response := httptest.NewRecorder()
	newServer(t).ServeHTTP(response, httptest.NewRequest(http.MethodGet, path, nil))
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
	response := post(newServer(t), "/games/demo/actions", `{"type":"End phase"}`)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	decoded := decodeView(t, response)
	if decoded.Turn != 1 || decoded.Phase != "Movement" {
		t.Errorf("view = turn %d, phase %q; want turn 1, phase \"Movement\"", decoded.Turn, decoded.Phase)
	}
}

func TestPostAnActionThatIsNotAvailableIsRejectedAndTheGameIsUnchanged(t *testing.T) {
	server := newServer(t)

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

func TestSimultaneousActionsAreAllApplied(t *testing.T) {
	const clicks = 100 // 20 whole turns of 5 phases
	server := newServer(t)

	var wg sync.WaitGroup
	for range clicks {
		wg.Add(1)
		go func() {
			defer wg.Done()
			post(server, "/games/demo/actions", `{"type":"End phase"}`)
		}()
	}
	wg.Wait()

	after := httptest.NewRecorder()
	server.ServeHTTP(after, httptest.NewRequest(http.MethodGet, "/games/demo", nil))
	decoded := decodeView(t, after)
	if decoded.Turn != 21 || decoded.Phase != "Upkeep" {
		t.Errorf("game = turn %d, phase %q; want turn 21, phase \"Upkeep\"", decoded.Turn, decoded.Phase)
	}
}

func TestPostWithAnOversizedBodyIsRejected(t *testing.T) {
	body := `{"type":"` + strings.Repeat("x", 1<<20) + `"}`

	response := post(newServer(t), "/games/demo/actions", body)

	if response.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("status = %d, want %d", response.Code, http.StatusRequestEntityTooLarge)
	}
}

func TestPostToAnUnknownGameReturnsNotFound(t *testing.T) {
	response := post(newServer(t), "/games/unknown/actions", `{"type":"End phase"}`)

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

// A game that refers to content that does not exist cannot be shown, so an
// Action on it fails and must leave the game as it was.
func TestPostThatCannotComposeTheViewFailsAndTheGameIsUnchanged(t *testing.T) {
	loaded, err := contentyaml.Load(os.DirFS("../../../content"))
	if err != nil {
		t.Fatalf("loading content: %v", err)
	}
	broken := engine.Game{Turn: 1, Phase: engine.PhaseUpkeep, AncientOne: "not-in-content"}
	store := games.NewMemory()
	store.Put("broken", broken)

	response := post(api.NewHandler(store, loaded), "/games/broken/actions", `{"type":"End phase"}`)

	if response.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", response.Code, http.StatusInternalServerError)
	}
	after, _ := store.Get("broken")
	if !reflect.DeepEqual(after, broken) {
		t.Errorf("game = %+v, want it unchanged: %+v", after, broken)
	}
}
