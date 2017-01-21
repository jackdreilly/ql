package quiklyrics

import (
	"errors"
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
		return Lyrics{}, errors.New("no title")
	}
	title := titleFind.Get(0).FirstChild.Data
	chords, e := find.Html()
	if e != nil {
		log.Println("Failed on html grab.", e.Error())
		return Lyrics{}, e
	}
	return Lyrics{chords, title}, nil
}

func GetChords(searchResults []Result) (Lyrics, []Alternative, error) {
	alts := []Alternative{}
	found := false
	sln := Lyrics{}
	for _, result := range searchResults {
		log.Println(result)
		var err error
		lyrics := Lyrics{}
		if strings.Contains(result.URL(), "ultimate") {
			lyrics, err = UltimateGuitar(result.URL())
		} else {
			continue
		}
		if err == nil {
			if !found {
				found = true
				sln = lyrics
				continue
			}
			alts = append(alts, Alternative{
				Title: strings.TrimSpace(lyrics.Title),
				Url:   result.URL(),
			})
		}
	}
	if !found {
		return Lyrics{}, alts, errors.New("No matches")
	}
	sln.Lyrics = strings.TrimSpace(sln.Lyrics)
	return sln, alts, nil

}

func GetChordsForQuery(query string) (Lyrics, []Alternative, error) {
	log.Println("Looking for chords", query)
	return GetChords(GoogleSearchChords(query))
}
