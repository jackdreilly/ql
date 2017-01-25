package quiklyrics

import (
	"bytes"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Lyrics struct {
	Lyrics string
	Title  string
}

func Genius(link string) (Lyrics, error) {
	doc, err := CurrentFetcher.Fetch(link)
	if err != nil {
		log.Println("could not fetch failed for " + link)
		return Lyrics{}, errors.New("no match")
	}
	find := doc.Find(".lyrics p")
	if find.Length() < 1 {
		log.Println("could not base parse failed for " + link)
		return Lyrics{}, errors.New("no match")
	}
	var buffer bytes.Buffer
	var child = find.Get(0).FirstChild
	for child != nil {
		for child.Data == "br" {
			child = child.NextSibling
			if child == nil {
				break
			}
		}
		if child == nil {
			break
		}
		for child.Data == "a" {
			buffer.WriteString(child.FirstChild.Data)
			child = child.NextSibling
			if child == nil {
				break
			}
		}
		if child == nil {
			break
		}
		if child.Data != "br" && child.Data != "a" {
			buffer.WriteString(child.Data)
			buffer.WriteString("\n")
		}
		child = child.NextSibling
	}
	titleFind := doc.Find("title")
	if titleFind.Length() < 1 {
		return Lyrics{}, errors.New("no title")
	}
	title := strings.Replace(strings.Replace(titleFind.Get(0).FirstChild.Data, " | LyricsFreak", "", 1), " Lyrics", "", 1)
	return Lyrics{buffer.String(), title}, nil
}

func LyricsFreak(link string) (Lyrics, error) {
	doc, err := CurrentFetcher.Fetch(link)
	if err != nil {
		log.Println("failed for " + link)
		return Lyrics{}, errors.New("no match")
	}
	find := doc.Find("#content_h")
	if find.Length() < 1 {
		log.Println("failed for " + link)
		return Lyrics{}, errors.New("no match")
	}
	var buffer bytes.Buffer
	var child = find.Get(0).FirstChild
	for child != nil {
		for child.Data == "br" {
			child = child.NextSibling
		}
		buffer.WriteString(child.Data)
		buffer.WriteString("\n")
		child = child.NextSibling
	}
	titleFind := doc.Find("title")
	if titleFind.Length() < 1 {
		return Lyrics{}, errors.New("no title")
	}
	title := strings.Replace(strings.Replace(titleFind.Get(0).FirstChild.Data, " | LyricsFreak", "", 1), " Lyrics", "", 1)
	return Lyrics{buffer.String(), title}, nil
}

func DirectLyrics(link string) (Lyrics, error) {
	doc, err := CurrentFetcher.Fetch(link)
	if err != nil {
		log.Println("failed for " + link)
		return Lyrics{}, errors.New("no match")
	}
	find := doc.Find(".lyrics p")
	if find.Length() < 1 {
		log.Println("failed for " + link)
		return Lyrics{}, errors.New("no match")
	}
	var buffer bytes.Buffer
	var child = find.Get(0).FirstChild
	for child != nil {
		for child.Data == "br" {
			child = child.NextSibling
		}
		buffer.WriteString(child.Data)
		child = child.NextSibling
	}
	titleFind := doc.Find("title")
	if titleFind.Length() < 1 {
		return Lyrics{}, errors.New("no title")
	}
	title := strings.Replace(titleFind.Get(0).FirstChild.Data, " LYRICS", "", 1)
	return Lyrics{buffer.String(), title}, nil
}

func AzLyrics(link string) (Lyrics, error) {
	doc, err := CurrentFetcher.Fetch(link)
	if err != nil {
		log.Println("no doc for " + link)
		return Lyrics{}, errors.New("no doc fetch")
	}
	find := doc.Find(".main-page .row .col-xs-12.col-lg-8.text-center div:not([class]):not([data-id])")
	if find.Length() < 1 {
		log.Println("no lyrics for " + link)
		return Lyrics{}, errors.New("no lyrics found")
	}
	var buffer bytes.Buffer
	var child = find.Get(0).FirstChild.NextSibling.NextSibling
	for child != nil {
		for child.Data == "br" {
			child = child.NextSibling
		}
		buffer.WriteString(child.Data)
		child = child.NextSibling
	}
	titleFind := doc.Find("title")
	if titleFind.Length() < 1 {
		log.Println("no lyrics for " + link)
		return Lyrics{}, errors.New("no title")
	}
	title := strings.Replace(titleFind.Get(0).FirstChild.Data, " LYRICS", "", 1)
	return Lyrics{buffer.String(), title}, nil
}

func MetroLyrics(link string) (Lyrics, error) {
	doc, err := CurrentFetcher.Fetch(link)
	if err != nil {
		log.Println("failed for " + link)
		return Lyrics{}, errors.New("no match")
	}
	titleFind := doc.Find("title")
	if titleFind.Length() < 1 {
		return Lyrics{}, errors.New("no title")
	}
	title := strings.Replace(strings.Replace(titleFind.Get(0).FirstChild.Data, " | MetroLyrics", "", 1), " Lyrics", "", 1)
	find := doc.Find(".verse")
	if find.Length() < 1 {
		return Lyrics{}, errors.New("no match")
	}
	var buffer bytes.Buffer
	find.Each(func(i int, s *goquery.Selection) {
		var child = s.Get(0).FirstChild
		for child != nil {
			for child.Data == "br" {
				child = child.NextSibling
			}
			buffer.WriteString(child.Data)
			child = child.NextSibling
		}
		buffer.WriteString("\n\n")
	})

	return Lyrics{buffer.String(), title}, nil
}

func GetLyricsUrl(url string) (lyrics Lyrics, err error) {
	if strings.Contains(url, "lyricsfreak.com") {
		lyrics, err = LyricsFreak(url)
	} else if strings.Contains(url, "metrolyrics.com") {
		lyrics, err = MetroLyrics(url)
	} else if strings.Contains(url, "directlyrics.com") {
		lyrics, err = DirectLyrics(url)
	} else if strings.Contains(url, "azlyrics.com") {
		lyrics, err = AzLyrics(url)
	} else if strings.Contains(url, "genius.com") {
		lyrics, err = Genius(url)
	} else {
		err = errors.New("No lyrics site match")
	}
	return
}

func helperLyrics(result Result, c chan lyricsError, i int, done chan struct{}) {
	log.Println("started lyrics", i)
	lyrics, err := GetLyricsUrl(result.URL())
	select {
	case c <- lyricsError{&lyrics, result.URL(), err, i}:
	case <-done:
		log.Println("We got told to leave early:", i)
		return
	}
	log.Println("finished lyrics", i)
}

func pullResultsFromChannel(results chan lyricsError, done chan struct{}, numRequests int) (Lyrics, []Alternative, error) {
	resultsSlice := make([]lyricsError, numRequests)
	for i := 0; i < numRequests; i++ {
		resultsSlice[i].er = errors.New("Never populated")
	}
	now := time.Now()
	found0 := false
	successes := 0
	for i := 0; i < numRequests; i++ {
		x := <-results
		if x.er != nil {
			resultsSlice[i].er = x.er
			continue
		}
		successes += 1
		resultsSlice[x.i] = x
		if x.i == 0 {
			found0 = true
		}
		elapsed := time.Now().Sub(now).Seconds()
		log.Println("elapsed time:", elapsed, "seconds")
		if successes >= 5 && elapsed >= 0.5 && found0 {
			log.Println("leaving early.")
			break
		}
	}
	sln := -1
	alts := []Alternative{}
	for i, le := range resultsSlice {
		if le.er != nil {
			continue
		}
		if sln < 0 {
			sln = i
		} else {
			alts = append(alts, Alternative{Title: le.lyrics.Title, Url: le.url})
		}
	}
	if sln < 0 {
		return Lyrics{}, alts, errors.New("No matches")
	}
	s := *resultsSlice[sln].lyrics

	s.Lyrics = strings.TrimSpace(s.Lyrics)
	return s, alts, nil
}

func GetLyrics(searchResults []Result) (Lyrics, []Alternative, error) {
	results := make(chan lyricsError)
	done := make(chan struct{})
	defer close(done)
	for i, r := range searchResults {
		go helperLyrics(r, results, i, done)
	}
	return pullResultsFromChannel(results, done, len(searchResults))
}

func GetLyricsForQuery(query string) (Lyrics, []Alternative, error) {
	return GetLyrics(GoogleSearch(query))
}
