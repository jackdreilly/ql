package quiklyrics

import "reillybrothers.net/jackdreilly/quiklyrics/redis"

var (
	Client = redis.NewQLClient()
)
