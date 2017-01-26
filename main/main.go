package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/jackdreilly/ql"
)

const (
	mainPage       = "static/quik_lyrics.html"
	notFoundPage   = "static/not_found.html"
	redirectPrefix = "/quik"
	staticDir      = "static"
	staticPrefix   = "/static/"
)

func HandleMain(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, mainPage)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, notFoundPage)
}

func HandleRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, redirectPrefix, http.StatusMovedPermanently)
}

func main() {
	log.SetOutput(os.Stdout)
	r := mux.NewRouter()
	r.PathPrefix(staticPrefix).Handler(http.StripPrefix(staticPrefix, http.FileServer(http.Dir(staticDir))))
	r.HandleFunc("/", HandleRedirect)
	r.HandleFunc(redirectPrefix, HandleMain)
	r.HandleFunc("/suggest", quiklyrics.SuggestServer)
	r.HandleFunc("/lyrics", quiklyrics.SearchServer)
	r.HandleFunc("/chords", quiklyrics.ChordsServer)
	r.HandleFunc("/url", quiklyrics.UrlServer)
	r.HandleFunc("/recent", quiklyrics.RecentServer)
	r.NotFoundHandler = http.HandlerFunc(notFound)
	http.Handle("/", r)
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	http.ListenAndServe(":"+port, nil)
}
