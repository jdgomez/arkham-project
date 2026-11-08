// Command arkham serves the game: the API under /api and the front end.
package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/jdgomez/arkham-project/backend/internal/api"
	"github.com/jdgomez/arkham-project/backend/internal/content"
	"github.com/jdgomez/arkham-project/backend/internal/games"
	"github.com/jdgomez/arkham-project/backend/internal/web"
)

func main() {
	addr := flag.String("addr", ":8080", "address to listen on")
	contentDir := flag.String("content", "content", "directory with the game content")
	flag.Parse()

	loaded, err := content.Load(os.DirFS(*contentDir))
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", api.NewHandler(games.NewDemo(loaded), loaded)))
	mux.Handle("/", web.Handler())

	listener, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}
	// With ":0" the system picks the port, so the log line is how callers learn it.
	log.Printf("listening on %s", listener.Addr())
	log.Fatal(http.Serve(listener, mux))
}
