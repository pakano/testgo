package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type Limiter struct {
	limiter  rate.Limiter
	lastTime time.Time
	key      string
}

func (l *Limiter) Allow() bool {
	return l.limiter.Allow()
}

type Limiters struct {
	limiters sync.Map
}

func NewLimiters() *Limiters {
	ls := &Limiters{}
	go ls.cleanLimiter()
	return ls
}

func (ls *Limiters) GetLimiter(key string, r rate.Limit, b int) *Limiter {
	if v, ok := ls.limiters.Load(key); ok {
		l := v.(*Limiter)
		l.lastTime = time.Now()
		return l
	}
	l := &Limiter{
		limiter:  *rate.NewLimiter(r, b),
		lastTime: time.Now(),
		key:      key,
	}
	ls.limiters.Store(key, l)
	return l
}

func (ls *Limiters) cleanLimiter() {
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	for {
		<-ticker.C
		ls.limiters.Range(func(key, value interface{}) bool {
			l := value.(*Limiter)
			if time.Now().Sub(l.lastTime) > time.Minute {
				ls.limiters.Delete(key)
			}
			return true
		})
	}
}

var limiters = NewLimiters()

func LimitMiddleWare(ctx *gin.Context) {
	ip := ctx.RemoteIP()
	l := limiters.GetLimiter(ip, rate.Limit(1), 1)
	if !l.Allow() {
		ctx.Status(http.StatusTooManyRequests)
		ctx.Abort()
	}
}
