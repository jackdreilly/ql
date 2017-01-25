package redis

import (
	"encoding/json"
	"time"

	"log"

	rds "gopkg.in/redis.v5"
)

const (
	redisDbIndex = 0
	keySearches  = "searches"
)

type DbClient struct {
	client *rds.Client
}

func NewQLClient() *DbClient {
	c := &DbClient{
		client: rds.NewClient(&rds.Options{
			Addr: "localhost:6379",
		}),
	}
	c.initialize()
	return c
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (d *DbClient) initialize() {
}

type Search struct {
	SearchTerm     string
	LyricsOrChords string
	Date           int64
	FormattedDate  string
}

func (d *DbClient) StoreSearch(s string, lOrC string) {
	date := time.Now().Unix()
	search := Search{
		SearchTerm:     s,
		Date:           date,
		LyricsOrChords: lOrC,
	}
	r, e := json.MarshalIndent(search, "", " ")
	check(e)
	log.Printf("store search: %v", r)
	check(d.client.LPush(keySearches, string(r)).Err())
}

func (d *DbClient) AllSearches(maxResults int) []Search {
	results := []Search{}
	lyrics := map[string]bool{
		"mattress that you stole":          true,
		"you stick around and it may show": true,
	}
	offsetStride := 10
	for {
		for offset := 0; ; offset += offsetStride {
			ssI := d.client.LRange(keySearches, int64(offset), int64(offset+maxResults))
			if ssI.Err() != nil {
				return results
			}
			r := ssI.Val()
			if len(r) == 0 {
				return results
			}
			result := Search{}
			for _, rr := range r {
				if json.Unmarshal([]byte(rr), &result) != nil {
					return results
				}
				_, b := lyrics[result.SearchTerm]
				if b {
					continue
				}
				lyrics[result.SearchTerm] = true
				result.FormattedDate = time.Unix(result.Date, 0).Format("March 9")
				results = append(results, result)
				if len(results) >= maxResults {
					return results
				}
			}
		}
	}
	return results
}
