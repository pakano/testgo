package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"test/util"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
)

var ctx = context.Background()
var rdb *redis.Client

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

	res, err := rdb.Set("pro-count", 50, 0).Result()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("set count:", res)
}

func seckill(c *gin.Context) {
	uid := "user-" + util.GetRanStr(8)

	//事务函数
	txf := func(tx *redis.Tx) error {
		n, err := tx.Get("pro-count").Int()
		if err != nil {
			return err
		}
		if n <= 0 {
			return errors.New("seckill out")
		}

		_, err = tx.TxPipelined(func(pipe redis.Pipeliner) error {
			_, err = pipe.Decr("pro-count").Result()
			if err != nil {
				return err
			}
			_, err = pipe.SAdd("users", uid).Result()
			if err != nil {
				return err
			}
			return nil
		})
		return err
	}

	//乐观锁
	if err := rdb.Watch(txf, "pro-count"); err != nil {
		if err == redis.TxFailedErr {
			c.JSON(202, err)
			return
		}
		c.JSON(201, err)
		return
	}
	c.JSON(200, "ok")
}

func main() {
	e := gin.Default()
	e.GET("/seckill", seckill)
	e.Run(":8080")
}
