package common

import "github.com/redis/go-redis/v9"

var Rdb *redis.Client

func InitRDB() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "139.9.103.236:6379",
		Password: "clock1234",
		DB:       0,
	})

	Rdb = rdb
}

func GetRDB() *redis.Client {
	return Rdb
}
