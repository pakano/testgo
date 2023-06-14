package main

import (
	"fmt"
	"test/util"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

type Person struct {
	gorm.Model
	Name string
}

const (
	DSNFormat = "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true"
)

var DB *gorm.DB

func InitMySQL() {
	//MYSQL
	connWrite := fmt.Sprintf(DSNFormat, "root", "123", "172.17.0.2", "3306", "test")
	connRead := fmt.Sprintf(DSNFormat, "root", "123", "172.17.0.3", "3306", "test")

	var ormLogger logger.Interface
	if gin.Mode() == "debug" {
		ormLogger = logger.Default.LogMode(logger.Info)
	} else {
		ormLogger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       connWrite, // DSN source
		DefaultStringSize:         256,       // string类型的字段长度
		DisableDatetimePrecision:  true,      // 禁止时间精度
		DontSupportRenameIndex:    true,      // 重命名索引时采用删除并重建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,      // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: true,      // 根据版本自动配置
	}), &gorm.Config{
		Logger: ormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Second * 30)
	DB = db

	//读写分离
	DB.Use(dbresolver.Register(
		dbresolver.Config{
			Sources:  []gorm.Dialector{mysql.Open(connWrite)},
			Replicas: []gorm.Dialector{mysql.Open(connRead)},
			Policy:   dbresolver.RandomPolicy{},
		},
	))

	// Migration()
	DB.AutoMigrate(&Person{})
}

func main() {
	InitMySQL()

	DB.Debug().Create(&Person{Name: util.GetRanStr(5)})
	var p Person
	DB.Debug().Model(&Person{}).First(&p)
	DB.Debug().Model(&Person{}).First(&p)
	DB.Debug().Model(&Person{}).First(&p)
	DB.Debug().Model(&Person{}).Unscoped().Where("id = ?", 10).Delete(&Person{})
}
