package db

import (
	"fmt"
	rdb "github.com/aurora/Filestore-server/db/redis"
	"testing"
)

func TestHGetAll(t *testing.T) {
	rConn := rdb.RedisConn()
	rConn.HMSet("MP_"+"123456789", map[string]interface{}{
		"chunkCount": 1,
		"fileHash":   "54546546532168sawrea8we",
		"FileName":   "test",
		"FileSize":   546431,
	})
	data := rConn.HGetAll("MP_" + "123456789").Val()
	for k, v := range data {
		fmt.Println(k, v)
	}
}
