package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Person struct {
	ID   int `gorm:"AUTO_INCREMENT"`
	Name string
	Age  int
}

var DSN = "root:123@tcp(mysql5740:3306)/sharing?charset=utf8mb4"

func InitMysql() {
	var err error
	DB, err = gorm.Open(mysql.Open(DSN), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	DB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").AutoMigrate(Person{})

	DB.Debug().Create(&Person{Name: "zhang3", Age: 10})
	DB.Debug().Create(&Person{Name: "li4", Age: 110})
}

var DB *gorm.DB
var RDB *redis.Client

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:         "redis1:6379",
		Password:     "", // no password set
		DB:           0,  // use default DB
		PoolSize:     8,
		MinIdleConns: 5,
	})

	_, err := RDB.Ping().Result()
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	//InitMysql()
	InitRedis()

	en := gin.Default()

	en.GET("/get/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		var p Person
		DB.Debug().First(&p, "id = ?", id)
		ctx.JSON(200, p)
	})

	en.GET("/redis/set", func(ctx *gin.Context) {
		key := ctx.Query("key")
		value := ctx.Query("value")

		_, err := RDB.Set(key, value, 0).Result()
		if err != nil {
			fmt.Println(err)
		}
		ctx.String(200, "ok")
	})

	en.GET("/redis/get", func(ctx *gin.Context) {
		key := ctx.Query("key")
		value, _ := RDB.Get(key).Result()
		ctx.JSON(200, gin.H{
			"key":   key,
			"value": value,
		})
	})
	en.Run(":80")
}
