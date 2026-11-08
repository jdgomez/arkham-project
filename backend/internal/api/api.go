// Package api is the HTTP adapter over the engine.
package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/jdgomez/arkham-project/backend/internal/engine"
	"github.com/jdgomez/arkham-project/backend/internal/games"
)

// maxBodyBytes caps a request body: an Action is a few bytes.
const maxBodyBytes = 1 << 10

// Games is where the API finds the games it serves. Update must be atomic and
// return games.ErrNotFound for an unknown id.
type Games interface {
	Get(id string) (engine.Game, bool)
	Update(id string, change func(engine.Game) (engine.Game, error)) (engine.Game, error)
}

// NewHandler returns the router of the API.
func NewHandler(store Games) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /games/{id}", func(w http.ResponseWriter, r *http.Request) {
		game, ok := store.Get(r.PathValue("id"))
		if !ok {
			http.NotFound(w, r)
			return
		}
		writeJSON(w, engine.ViewOf(game))
	})
	mux.HandleFunc("POST /games/{id}/actions", func(w http.ResponseWriter, r *http.Request) {
		var action engine.Action
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBodyBytes)).Decode(&action); err != nil {
			var tooLarge *http.MaxBytesError
			if errors.As(err, &tooLarge) {
				http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
				return
			}
			http.Error(w, "invalid action", http.StatusBadRequest)
			return
		}
		next, err := store.Update(r.PathValue("id"), func(game engine.Game) (engine.Game, error) {
			return engine.Apply(game, action)
		})
		switch {
		case errors.Is(err, games.ErrNotFound):
			http.NotFound(w, r)
		case errors.Is(err, engine.ErrActionNotAvailable):
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		case err != nil:
			log.Printf("applying action: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		default:
			writeJSON(w, engine.ViewOf(next))
		}
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
