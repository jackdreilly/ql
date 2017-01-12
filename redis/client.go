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
	SearchTerm string
	Date       int64
}

type Website struct {
	Url String
}

func (d *DbClient) StoreSearch(s string) {
	date := time.Now().Unix()
	search := Search{
		SearchTerm: s,
		Date:       date,
	}
	r, e := json.MarshalIndent(search, "", " ")
	check(e)
	log.Printf("store search: %v", r)
	check(d.client.LPush(keySearches, string(r)).Err())
}

func (d *DbClient) AllSearches(maxResults int64) []Search {
	ssI := d.client.LRange(keySearches, 0, maxResults)
	check(ssI.Err())
	ss := ssI.Val()
	r := make([]Search, len(ss))
	for i, s := range ss {
		check(json.Unmarshal([]byte(s), &r[i]))
	}
	return r
}
