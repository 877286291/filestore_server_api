package redis

import (
	"github.com/go-redis/redis"
	"log"
)

var rdb *redis.Client

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	_, err := rdb.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}
}
func RedisConn() *redis.Client {
	return rdb
}
