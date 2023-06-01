package middleware

import (
	"log"
	"os"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
)

func MyRecover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if e := recover(); e != nil {
				pwd, _ := os.Getwd()
				fd, e2 := os.Create(pwd + "/" + time.Now().Format("20060102030405") + ".stack")
				if e2 == nil && fd != nil {
					e3, ok := e.(error)
					if ok {
						fd.WriteString("recover:" + e3.(error).Error() + "\n")
					}
					fd.Write(debug.Stack())
					fd.Close()
				} else {
					log.Println("error:", e, e2)
					log.Println("error:", string(debug.Stack()))
				}
			}
		}()
		c.Next() // 调用下一个处理
	}
}
