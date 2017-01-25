package quiklyrics

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

func UltimateGuitar(link string) (Lyrics, error) {
	doc, err := CurrentFetcher.Fetch(link)
	if err != nil {
		log.Println("failed for " + link)
		return Lyrics{}, errors.New("no match")
	}
	find := doc.Find("pre.js-tab-content")
	if find.Length() < 1 {
		log.Println("failed for " + link)
		return Lyrics{}, errors.New("no match")
	}
	titleFind := doc.Find("h1")
	if titleFind.Length() < 1 {
		log.Println("Cannot find title.")
		return Lyrics{}, errors.New("no song title")
	}
	songTitle := titleFind.Get(0).FirstChild.Data
	titleFind = doc.Find(".t_autor a")
	if titleFind.Length() < 1 {
		log.Println("Cannot find title.")
		return Lyrics{}, errors.New("no author")
	}
	songAuthor := titleFind.Get(0).FirstChild.Data

	chords, e := find.Html()
	if e != nil {
		log.Println("Failed on html grab.", e.Error())
		return Lyrics{}, e
	}
	return Lyrics{chords, fmt.Sprintf("%v - %v", songAuthor, songTitle)}, nil
}

type lyricsError struct {
	lyrics *Lyrics
	url    string
	er     error
	i      int
}

func GetChordsUrl(url string) (lyrics Lyrics, err error) {
	if strings.Contains(url, "ultimate-guitar.com") {
		lyrics, err = UltimateGuitar(url)
	} else {
		err = errors.New("No chords site match")
	}
	return
}

func helper(result Result, c chan lyricsError, i int, done chan struct{}) {
	log.Println("started chords", i)
	lyrics, err := GetChordsUrl(result.URL())
	select {
	case c <- lyricsError{&lyrics, result.URL(), err, i}:
	case <-done:
		log.Println("We got told to leave early:", i)
		return
	}
	log.Println("finished lyrics", i)
}

func GetChords(searchResults []Result) (Lyrics, []Alternative, error) {
	results := make(chan lyricsError, 10)
	done := make(chan struct{})
	defer close(done)
	for i, r := range searchResults {
		go helper(r, results, i, done)
	}
	return pullResultsFromChannel(results, done, len(searchResults))
}

func GetChordsForQuery(query string) (Lyrics, []Alternative, error) {
	return GetChords(GoogleSearchChords(query))
}
