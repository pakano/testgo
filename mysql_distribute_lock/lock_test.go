package main

import (
	"fmt"
	"os"
	"sync"
	"test/util"
	"testing"
	"time"

	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	dsn := "root:123@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = util.NewMysqlInstance(dsn)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	db.AutoMigrate(Locker{})
}

func Test001(t *testing.T) {
	lockKey := "key1"
	lock := NewRMutex(db)
	go func() {
		requestID := GetRequestID()
		for i := 0; i < 1; i++ {
			lock.Lock(lockKey, requestID, 10000, 1000)
		}
		time.Sleep(time.Second * 5)
		for i := 0; i < 0; i++ {
			lock.Unlock(lockKey, requestID)
		}
	}()
	time.Sleep(time.Second * 60)
	go func() {
		requestID := GetRequestID()
		for i := 0; i < 1; i++ {
			lock.Lock(lockKey, requestID, 1000, 10000)
		}

		for i := 0; i < 1; i++ {
			lock.Unlock(lockKey, requestID)
		}
	}()
	select {}
}

func Test002(t *testing.T) {
	var num int32
	lockKey := "key2"
	wg := sync.WaitGroup{}
	lock := NewRMutex(db)

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			requestID := GetRequestID()
			lock.Lock(lockKey, requestID, 10000, 100000)
			defer lock.Unlock(lockKey, requestID)

			num++
		}()
	}
	wg.Wait()
	fmt.Println(num)
}
