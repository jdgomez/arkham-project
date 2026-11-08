// Package api is the HTTP adapter over the engine.
package api

import (
	"encoding/json"
	"net/http"

	"github.com/jdgomez/arkham-project/backend/internal/engine"
)

// Games is where the API finds the games it serves.
type Games interface {
	Get(id string) (engine.Game, bool)
}

// NewHandler returns the router of the API.
func NewHandler(games Games) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /games/{id}", func(w http.ResponseWriter, r *http.Request) {
		game, ok := games.Get(r.PathValue("id"))
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(engine.ViewOf(game))
	})
	return mux
}
