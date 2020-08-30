package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	baseURL = "http://suggestqueries.google.com/complete/search?output=toolbar&hl=en&q="
)

// url gives suggest url for query string
func suggestUrl(query string, lyricsOrChords string) string {
	lyricsQuery := fmt.Sprintf("%v %v", lyricsOrChords, query)
	return baseURL + url.QueryEscape(lyricsQuery)
}

// Suggest uses google suggest to return suggestions for current query
func Suggest(query string, lyricsOrChords string) []string {
	doc, err := CurrentFetcher.Fetch(suggestUrl(query, lyricsOrChords))
	if err != nil {
		panic(err)
	}
	var results = make([]string, 0)
	doc.Find("suggestion").Each(func(i int, s *goquery.Selection) {
		data := s.AttrOr("data", "")
		if len(data) > 0 {
			results = append(results, strings.Replace(data, lyricsOrChords+" ", "", 1))
		}
	})
	return results
}
