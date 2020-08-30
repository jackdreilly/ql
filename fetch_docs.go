package main

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type DocFetcher interface {
	Fetch(urlString string) (*goquery.Document, error)
	Get(urlString string) (*http.Response, error)
}

type standardFetcher int

func (s standardFetcher) Fetch(urlString string) (*goquery.Document, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return goquery.NewDocumentFromReader(resp.Body)
}

func (s standardFetcher) Get(urlString string) (*http.Response, error) {
	return http.Get(urlString)
}

const (
	Fetcher = standardFetcher(1)
)

var (
	CurrentFetcher DocFetcher = Fetcher
)
