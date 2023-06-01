package main

import (
	"fmt"
	"os"
	"sync"
	"test/util"
	"time"

	godisson "github.com/cheerego/go-redisson"
	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client
var g *godisson.Godisson

func init() {
	var err error
	rdb, err = util.NewRedisInstance()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	g = godisson.NewGodisson(rdb, godisson.WithWatchDogTimeout(30*time.Second))
	if g == nil {
		fmt.Println("g is nil")
		os.Exit(-1)
	}
}

func test01() {
	key := "key1"
	lock := g.NewMutex(key)
	var num int
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			lock.TryLock(int64(time.Millisecond*10000), int64(time.Millisecond*10000))
			defer lock.Unlock()
			num++
		}()
	}
	wg.Wait()

	fmt.Println(num)
}

func test02() {
	key := "key1"
	lock := g.NewMutex(key)

	lock.TryLock(10000, -1)
	//defer lock.Unlock()
}

func main() {
	start := time.Now()
	test01()
	fmt.Println(time.Now().Sub(start))
}
