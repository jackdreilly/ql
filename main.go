package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	log.SetOutput(os.Stdout)
	r := mux.NewRouter()
	r.HandleFunc("/suggest", SuggestServer)
	r.HandleFunc("/lyrics", SearchServer)
	http.Handle("/", handlers.CORS()(r))
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	http.ListenAndServe(":"+port, nil)
}
