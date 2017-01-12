package main

import (
	"log"
	"net/http"
	"os"
	"reillybrothers.net/jackdreilly/quiklyrics"
)

func main() {
	log.SetOutput(os.Stdout)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	http.HandleFunc("/suggest", quiklyrics.SuggestServer)
	http.HandleFunc("/lyrics", quiklyrics.SearchServer)
	http.ListenAndServe(":8083", nil)
}
