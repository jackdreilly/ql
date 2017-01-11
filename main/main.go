package main

import (
	"github.com/jackdreilly/quiklyrics"
	"log"
	"net/http"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	http.HandleFunc("/suggest", quiklyrics.SuggestServer)
	http.HandleFunc("/lyrics", quiklyrics.SearchServer)
	http.ListenAndServe(":8083", nil)
}
