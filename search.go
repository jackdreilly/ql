package quiklyrics

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/url"

	"google.golang.org/api/customsearch/v1"
)

const (
	key               = "AIzaSyDTtyBAmPzhyB2vuHfUTP7k5Wbvv24PGsg"
	lyricsEngine      = "001283733761018543620%3Aujg1bejtz_y"
	lyricsEngineNoAz  = "001283733761018543620:q_vn9d9funi"
	chordsEngine      = "001283733761018543620:wwna3zpolnc"
	lyricsBaseURL     = "https://www.googleapis.com/customsearch/v1?cx=" + lyricsEngine + "&key=" + key + "&q="
	lyricsNoAzBaseURL = "https://www.googleapis.com/customsearch/v1?cx=" + lyricsEngineNoAz + "&key=" + key + "&q="
	chordsBaseURL     = "https://www.googleapis.com/customsearch/v1?cx=" + chordsEngine + "&key=" + key + "&q="
)

type Result struct {
	resultUrl string
	title     string
}

func GoogleSearch(query string) []Result {
	requestURL := lyricsBaseURL + url.QueryEscape("lyrics "+query)
	res, err := CurrentFetcher.Get(requestURL)
	if err != nil {
		log.Println(requestURL)
		log.Println(query)
		panic(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var data customsearch.Search
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err)
	}
	var results = make([]Result, 0)
	for _, item := range data.Items {
		results = append(results, Result{item.Link, item.Title})
	}
	return results
}

func GoogleSearchNoAz(query string) []Result {
	requestURL := lyricsNoAzBaseURL + url.QueryEscape("lyrics "+query)
	res, err := CurrentFetcher.Get(requestURL)
	if err != nil {
		log.Println(requestURL)
		log.Println(query)
		panic(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var data customsearch.Search
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err)
	}
	var results = make([]Result, 0)
	for _, item := range data.Items {
		results = append(results, Result{item.Link, item.Title})
	}
	return results
}

func GoogleSearchChords(query string) []Result {
	requestURL := chordsBaseURL + url.QueryEscape("chords "+query)
	res, err := CurrentFetcher.Get(requestURL)
	if err != nil {
		log.Println(requestURL)
		log.Println(query)
		panic(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var data customsearch.Search
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err)
	}
	var results = make([]Result, 0)
	for _, item := range data.Items {
		results = append(results, Result{item.Link, item.Title})
	}
	return results
}

func (r *Result) URL() string {
	return r.resultUrl
}
