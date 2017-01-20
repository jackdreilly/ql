package quiklyrics

import (
	"encoding/json"
	"net/http"
)

type SuggestResponse struct {
	Suggestions []string `json:"suggestions"`
}

type Alternative struct {
	Url   string
	Title string
}

type SearchResponse struct {
	Lyrics       Lyrics
	Alternatives []Alternative
}

func suggestFake(query string) []string {
	return []string{query, query, query + " extra", "Hey ya", "Boo yeah"}
}

func searchFake(query string) (Lyrics, []Alternative) {
	return Lyrics{"fake lyrics " + query, "Fake Title" + query}, []Alternative{
		Alternative{
			Url:   "urlA-" + query,
			Title: "titleA-" + query,
		},
		Alternative{
			Url:   "urlB-" + query,
			Title: "titleB-" + query,
		},
	}
}

func SuggestServer(w http.ResponseWriter, r *http.Request) {
	response := Suggest(r.URL.Query().Get("lyrics"))
	// response := suggestFake(r.URL.Query().Get("lyrics"))
	b, e := json.Marshal(SuggestResponse{response})
	if e != nil {
		panic(e)
	}
	w.Write(b)
}

func SearchServer(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("lyrics")
	Client.StoreSearch(query)
	lyrics, alts, e := GetLyricsForQuery(query)
	if e != nil {
		http.Error(w, e.Error(), http.StatusNotFound)
		return
	}
	// lyrics, alts := searchFake(query)
	b, e := json.Marshal(SearchResponse{
		Lyrics:       lyrics,
		Alternatives: alts,
	})
	if e != nil {
		panic(e)
	}
	w.Write(b)
}

var (
	ChordsServer = SearchServer
)
