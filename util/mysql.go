package util

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// dsn := "root:123@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
// dsn: user:passwd@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
func NewMysqlInstance(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(1000)
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour * 1)
	sqlDB.SetConnMaxIdleTime(time.Hour * 24)

	return db, nil
}
