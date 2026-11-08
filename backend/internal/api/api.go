// Package api is the HTTP adapter over the engine.
package api

import (
	"bytes"
	"encoding/json"
	"log"
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
		writeJSON(w, engine.ViewOf(game))
	})
	return mux
}

// writeJSON encodes the whole body before sending anything, so a failure
// becomes a 500 and never a 200 with half a body.
func writeJSON(w http.ResponseWriter, body any) {
	var encoded bytes.Buffer
	if err := json.NewEncoder(&encoded).Encode(body); err != nil {
		log.Printf("encoding response: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(encoded.Bytes())
}
