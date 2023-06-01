package main

import (
	"strconv"
	"test/jwt/hs256"

	"github.com/gin-gonic/gin"
)

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(ctx *gin.Context) {
	var req LoginReq
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.String(200, "param error")
		return
	}
	if req.Username != "zhangsan" || req.Password != "123" {
		ctx.String(200, "auth error")
		return
	}
	user := hs256.User{UserID: 1, Username: req.Username, GrantScope: "xx"}
	token, err := hs256.GenerateTokenUsingHs256(user)
	if err != nil {
		ctx.String(200, "unknown error")
		return
	}
	ctx.JSON(200, gin.H{
		"id":       1,
		"username": req.Username,
		"token":    token,
	})
}

func UserInfo(ctx *gin.Context) {
	userid := ctx.Param("userid")
	username := ctx.Param("username")
	ctx.JSON(200, gin.H{
		"userid":   userid,
		"username": username,
	})
}

func Auth(ctx *gin.Context) {
	token := ctx.Request.Header.Get("Authorization")
	if len(token) == 0 {
		ctx.String(401, "No Authorization")
		ctx.Abort()
	}
	claim, err := hs256.ParseTokenHs256(token)
	if err != nil {
		ctx.String(401, "No Authorization")
		ctx.Abort()
	}
	ctx.Params = append(ctx.Params, gin.Param{
		Key:   "userid",
		Value: strconv.Itoa(claim.User.UserID),
	})
	ctx.Params = append(ctx.Params, gin.Param{
		Key:   "username",
		Value: claim.User.Username,
	})
	ctx.Next()
}

func main() {
	en := gin.Default()
	en.POST("/login", Login)
	g := en.Group("/user", Auth)
	{
		g.GET("/userinfo", UserInfo)
	}
	en.Run(":8081")
}
