// Command arkham serves the game: the API under /api and the front end.
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/jdgomez/arkham-project/backend/internal/api"
	"github.com/jdgomez/arkham-project/backend/internal/games"
	"github.com/jdgomez/arkham-project/backend/internal/web"
)

func main() {
	addr := flag.String("addr", ":8080", "address to listen on")
	flag.Parse()

	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", api.NewHandler(games.NewDemo())))
	mux.Handle("/", web.Handler())

	log.Printf("listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, mux))
}
