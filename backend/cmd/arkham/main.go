// Command arkham serves the game: the API under /api and the front end.
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/jdgomez/arkham-project/backend/internal/api"
	"github.com/jdgomez/arkham-project/backend/internal/engine"
	"github.com/jdgomez/arkham-project/backend/internal/games"
	"github.com/jdgomez/arkham-project/backend/internal/web"
)

const demoGameID = "demo"

func main() {
	addr := flag.String("addr", ":8080", "address to listen on")
	flag.Parse()

	store := games.NewMemory()
	store.Put(demoGameID, engine.DemoGame())

	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", api.NewHandler(store)))
	mux.Handle("/", web.Handler())

	log.Printf("listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, mux))
}
