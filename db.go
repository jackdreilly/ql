package quiklyrics

import "github.com/jackdreilly/ql/redis"

var (
	Client = redis.NewQLClient()
)
