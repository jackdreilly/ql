package quiklyrics

import (
	"encoding/json"
	"log"
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
	response := Suggest(r.URL.Query().Get("lyrics"), r.URL.Query().Get("lyricsOrChords"))
	// response := suggestFake(r.URL.Query().Get("lyrics"))
	b, e := json.Marshal(SuggestResponse{response})
	if e != nil {
		panic(e)
	}
	w.Write(b)
}

func SearchServer(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("lyrics")
	lOrC := r.URL.Query().Get("lyricsOrChords")
	Client.StoreSearch(query, lOrC)
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

func UrlServer(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	// TODO: Store redis.
	lyricsOrChords := r.URL.Query().Get("lyricsOrChords")
	log.Println("url", url, "lyrics or chords", lyricsOrChords)
	var lyrics Lyrics
	var err error
	switch lyricsOrChords {
	case "lyrics":
		log.Println("Get lyrics url", url)
		lyrics, err = GetLyricsUrl(url)
	case "chords":
		log.Println("Get chords url", url)
		lyrics, err = GetChordsUrl(url)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	b, err := json.Marshal(lyrics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Write(b)
}

func ChordsServer(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("lyrics")
	lOrC := r.URL.Query().Get("lyricsOrChords")
	Client.StoreSearch(query, lOrC)
	lyrics, alts, e := GetChordsForQuery(query)
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
		http.Error(w, e.Error(), http.StatusNotFound)
		return
	}
	w.Write(b)
}

func RecentServer(w http.ResponseWriter, r *http.Request) {
	searches := Client.AllSearches(10)
	b, e := json.Marshal(searches)
	if e != nil {
		http.Error(w, e.Error(), http.StatusNotFound)
		return
	}
	w.Write(b)
}
