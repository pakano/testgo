package delaydb

import (
	"context"
	"log"
	"time"

	"gorm.io/gorm"
)

type dalayedData struct {
	args []interface{}
	fn   func(db *gorm.DB, data ...interface{}) error
}

type DelayedQueue struct {
	db       *gorm.DB
	maxSize  int
	queue    []dalayedData
	duration time.Duration
	timer    *time.Timer
	chanData chan dalayedData
	chanTime chan struct{}
}

func NewDelayedQueue(db *gorm.DB, maxSize int, duration time.Duration) *DelayedQueue {
	return &DelayedQueue{
		db:       db,
		maxSize:  maxSize,
		queue:    make([]dalayedData, 0),
		chanData: make(chan dalayedData, 1000),
		chanTime: make(chan struct{}),
	}
}

func (queue *DelayedQueue) Run(ctx context.Context) {
	log.Println("info:", "delayed queue running...")
	for {
		select {
		case data := <-queue.chanData:
			queue.queue = append(queue.queue, data)
			if len(queue.queue) > queue.maxSize {
				tx := queue.db.Begin()
				for _, d := range queue.queue {
					d.fn(tx, d.args...)
				}
				tx.Debug().Commit()
				queue.queue = make([]dalayedData, 0)
			} else {
				if queue.timer == nil {
					queue.timer = time.AfterFunc(queue.duration, func() {
						queue.chanTime <- struct{}{}
					})
				}
			}

		case <-queue.chanTime:
			if len(queue.queue) > 0 {
				tx := queue.db.Begin()
				for _, d := range queue.queue {
					d.fn(tx, d.args...)
				}
				tx.Commit()
			}
			queue.queue = make([]dalayedData, 0)
			queue.timer = nil

		case <-ctx.Done():
			log.Println("err:", "context canceld")
			return
		}
	}
}

func (queue *DelayedQueue) Insert(fn func(db *gorm.DB, args ...interface{}) (err error), data ...interface{}) {
	var args []interface{}
	args = append(args, data...)
	queue.chanData <- dalayedData{args: args, fn: fn}
}
