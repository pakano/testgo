package main

import (
	"context"
	"fmt"
	"log"
	"test/util"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
)

var seckillScript = `
local prodid=KEYS[1];
local userid=KEYS[2];
local userkey=KEYS[3];
local userExists=redis.call("sismember",userkey,userid);
if tonumber(userExists)==1 then
	return 2;
end
local num=redis.call("get",prodid);
if tonumber(num)<=0 then
	return 1;
else
	redis.call("decr",prodid);
	redis.call("SAdd",userkey,userid);
end
return 0;
`

var ctx = context.Background()
var rdb *redis.Client
var script *redis.Script

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:         "127.0.0.1:6379",
		Password:     "", // no password set
		DB:           0,  // use default DB
		PoolSize:     8,
		MinIdleConns: 5,
	})

	_, err := rdb.Ping().Result()
	if err != nil {
		log.Fatalln(err)
	}

	rdb.FlushAll()

	script = redis.NewScript(seckillScript)
	if script == nil {
		log.Fatalln("load script failed")
	}

	res, err := rdb.Set("pro-count", 50, 0).Result()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("set count:", res)
}

func seckill(c *gin.Context) {
	uid := "user-" + util.GetRanStr(8)

	v, err := script.Run(rdb, []string{"pro-count", uid, "users"}).Result()
	if err != nil {
		log.Fatalln(err)
	}
	n := v.(int64)
	switch n {
	case 0:
		c.JSON(200, "ok")
	case 1:
		c.JSON(201, "has seckill")
	case 2:
		c.JSON(202, "seckill out")
	}
}

func main() {
	e := gin.Default()
	e.GET("/seckill", seckill)
	e.Run(":8080")
}
