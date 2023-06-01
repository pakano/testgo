package main

import (
	"context"
	"fmt"
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

type RMutex struct {
	db     *gorm.DB
	cancel context.CancelFunc
	ctx    context.Context
}

func NewRMutex(db *gorm.DB) *RMutex {
	if db == nil {
		return nil
	}
	return &RMutex{db: db}
}

func GetRequestID() string {
	u := uuid.New()
	return u.String()
}

func (rm *RMutex) Lock(lockKey, requestID string, lease, timeout int64) bool {
	lockResult := false
	startTime := time.Now().UnixMilli()
	for {
		var lock Locker
		err := rm.db.Debug().First(&lock, "lock_key = ?", lockKey).Error
		if err != nil {
			//插入一条记录，被其他结点插入也无所谓
			rm.db.Debug().Create(&Locker{LockKey: lockKey, RequestID: "", LockCount: 0, Timeout: 0, Version: 0})
		} else {
			requestID2 := lock.RequestID
			if len(requestID2) == 0 {
				lock.RequestID = requestID
				lock.LockCount = 1
				lock.Timeout = time.Now().UnixMilli() + lease
				if rm.update(lock) == 1 {
					rm.renewExpirationScheduler(lockKey, requestID)
					lockResult = true
					break
				}
			} else if requestID2 == requestID {
				lock.LockCount = lock.LockCount + 1
				lock.Timeout = time.Now().UnixMilli() + lease
				if rm.update(lock) == 1 {
					lockResult = true
					break
				}
			} else {
				if lock.Timeout < time.Now().UnixMilli() {
					rm.resetLock(lock)
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

func (rm *RMutex) Unlock(lockKey, requestID string) {
	for {
		var lock Locker
		err := rm.db.Debug().First(&lock, "lock_key = ?", lockKey).Error
		if err == nil && lock.RequestID == requestID && lock.LockCount > 0 {
			if lock.LockCount == 1 {
				if rm.resetLock(lock) == 1 {
					rm.cancelExpirationRenewal()
					break
				}
			} else {
				lock.LockCount--
				if rm.update(lock) == 1 {
					break
				}
			}
		} else {
			break
		}
	}
}

func (rm *RMutex) resetLock(lock Locker) int64 {
	lock.RequestID = ""
	lock.LockCount = 0
	lock.Timeout = 0
	return rm.update(lock)
}

func (rm *RMutex) update(lock Locker) int64 {
	return rm.db.Debug().Model(Locker{}).Where("lock_key = ? AND version = ?", lock.LockKey, lock.Version).Updates(
		map[string]interface{}{
			"request_id": lock.RequestID,
			"lock_count": lock.LockCount,
			"timeout":    lock.Timeout,
			"version":    lock.Version + 1,
		},
	).RowsAffected
}

func (rm *RMutex) renewExpirationScheduler(lockKey, requestID string) {
	ctx, cancel = context.WithCancel(context.TODO())
	go rm.renewExpirationSchedulerGoroutine(lockKey, requestID)
}

var watchDogTimeout = time.Second * 30
var ctx context.Context
var cancel context.CancelFunc

func (rm *RMutex) renewExpirationSchedulerGoroutine(lockKey, requestID string) {
	ticker := time.NewTicker(watchDogTimeout / 3)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ret := rm.renewExpiration(lockKey, requestID)
			if ret == 0 {
				rm.cancelExpirationRenewal()
				fmt.Println("exit1")
				return
			}

		case <-ctx.Done():
			fmt.Println("exit2")
			return
		}
	}
}

func (rm *RMutex) renewExpiration(lockKey, requestID string) int64 {
	var ok bool
	for {
		var lock Locker
		err := rm.db.Debug().First(&lock, "lock_key = ?", lockKey).Error
		if err != nil {
			break
		} else if lock.RequestID == requestID && lock.LockCount > 0 {
			lock.Timeout = time.Now().Add(watchDogTimeout).UnixMilli()
			if rm.update(lock) == 1 {
				ok = true
				break
			}
		} else {
			break
		}
	}

	if ok {
		return 1
	}
	return 0
}

func (rm *RMutex) cancelExpirationRenewal() {
	if rm.cancel != nil {
		rm.cancel()
		rm.cancel = nil
	}
}
