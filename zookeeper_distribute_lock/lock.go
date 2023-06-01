package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

func main() {
	conn, _, err := zk.Connect([]string{"127.0.0.1:2181"}, 5*time.Second)
	if err != nil {
		panic(err)
	}

	lock := zk.NewLock(conn, "/lock", zk.WorldACL(zk.PermAll))

	var num int
	wg := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			lock.Lock()
			num++
			fmt.Println(num)
			lock.Unlock()
		}()
	}
	wg.Wait()
	fmt.Print(num)
}
