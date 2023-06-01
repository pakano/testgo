package main

import (
	"fmt"
	"os"
	"sync"
	"test/util"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Locker  '锁信息表'
type Locker struct {
	LockKey   string `gorm:"primaryKey"` //'锁唯一标志'
	RequestID string //'用来标识请求对象的'
	LockCount int    //'当前上锁次数'
	Timeout   int64  //'锁超时时间'
	Version   int    //'版本号，每次更新+1'
}

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

func GetRequestID() string {
	u := uuid.New()
	return u.String()
}

func Lock(lockKey, requestID string, lease, timeout int64) bool {
	lockResult := false
	startTime := time.Now().UnixMilli()
	for {
		var lock Locker
		err := db.Debug().First(&lock, "lock_key = ?", lockKey).Error
		if err != nil {
			//插入一条记录，被其他结点插入也无所谓
			db.Debug().Create(&Locker{LockKey: lockKey, RequestID: "", LockCount: 0, Timeout: 0, Version: 0})
		} else {
			requestID2 := lock.RequestID
			if len(requestID2) == 0 {
				lock.RequestID = requestID
				lock.LockCount = 1
				lock.Timeout = time.Now().UnixMilli() + lease
				if update(lock) == 1 {
					lockResult = true
					break
				}
			} else if requestID2 == requestID {
				lock.LockCount = lock.LockCount + 1
				lock.Timeout = time.Now().UnixMilli() + lease
				if update(lock) == 1 {
					lockResult = true
					break
				}
			} else {
				if lock.Timeout < time.Now().UnixMilli() {
					resetLock(lock)
				} else {
					if timeout+startTime < time.Now().UnixMilli() {
						break
					} else {
						time.Sleep(time.Millisecond * 100)
					}
				}
			}
		}
	}
	return lockResult
}

func Unlock(lockKey, requestID string) {
	var lock Locker
	err := db.Debug().First(&lock, "lock_key = ?", lockKey).Error
	if err == nil && lock.RequestID == requestID && lock.LockCount > 0 {
		if lock.LockCount == 1 {
			resetLock(lock)
		} else {
			lock.LockCount--
			update(lock)
		}
	}
}

func resetLock(lock Locker) int64 {
	lock.RequestID = ""
	lock.LockCount = 0
	lock.Timeout = 0
	return update(lock)
}

func update(lock Locker) int64 {
	return db.Debug().Model(Locker{}).Where("lock_key = ? AND version = ?", lock.LockKey, lock.Version).Updates(
		map[string]interface{}{
			"request_id": lock.RequestID,
			"lock_count": lock.LockCount,
			"timeout":    lock.Timeout,
			"version":    lock.Version + 1,
		},
	).RowsAffected
}

func test01() {
	lockKey := "key1"
	go func() {
		requestID := GetRequestID()
		for i := 0; i < 1; i++ {
			Lock(lockKey, requestID, 10000, 1000)
		}
		time.Sleep(time.Second * 5)
		for i := 0; i < 1; i++ {
			Unlock(lockKey, requestID)
		}
	}()
	time.Sleep(time.Second)
	go func() {
		requestID := GetRequestID()
		for i := 0; i < 1; i++ {
			Lock(lockKey, requestID, 1000, 10000)
		}

		for i := 0; i < 1; i++ {
			Unlock(lockKey, requestID)
		}
	}()
	select {}
}

func test02() {
	var num int32
	lockKey := "key2"
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			requestID := GetRequestID()
			Lock(lockKey, requestID, 10000, 100000)
			defer Unlock(lockKey, requestID)

			num++
		}()
	}
	wg.Wait()
	fmt.Println(num)
}

func main() {
	start := time.Now()
	test02()
	fmt.Println(time.Now().Sub(start))
}
