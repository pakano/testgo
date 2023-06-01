package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func server1() {
	en := gin.New()
	en.GET("/file", func(ctx *gin.Context) {
		ctx.String(200, "file1")
	})
	en.Run(":8001")
}

func server2() {
	en := gin.New()
	en.GET("/file", func(ctx *gin.Context) {
		fmt.Println(ctx.Request.RemoteAddr)
		fmt.Println(ctx.Request.Host)
		fmt.Println(ctx.GetHeader("x-real-ip"))
		ctx.String(200, "file2")
	})
	en.Run(":8002")
}

func main() {
	go server1()
	go server2()
	select {}
}
