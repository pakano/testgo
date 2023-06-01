package main

import (
	"errors"
	"fmt"
	"test/util"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Goods struct {
	ID      int `gorm:"primaryKey"`
	Name    string
	Count   int
	Version int
}

func main() {
	dsn := "root:123@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := util.NewMysqlInstance(dsn)
	if err != nil {
		fmt.Println(err)
		return
	}
	db.Exec("drop table if exists goods")
	db.AutoMigrate(Goods{})
	db.Debug().Create(&Goods{ID: 1, Name: "apple", Count: 2, Version: 0})

	for i := 0; i < 100; i++ {
		go func() {
			err := seckill_opt(db)
			fmt.Println(err)
		}()
	}
	time.Sleep(time.Second * 3)
}

func seckill(db *gorm.DB) error {
	// 再唠叨一下，事务一旦开始，你就应该使用 tx 处理数据
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	var good Goods
	if err := tx.Debug().Clauses(clause.Locking{
		Strength: "UPDATE",
	}).First(&good, "name = ?", "apple").Error; err != nil {
		tx.Rollback()
		return err
	}

	if good.Count <= 0 {
		tx.Rollback()
		return errors.New("no goods")
	}

	if tx.Debug().Model(Goods{}).Where("name = ?", good.Name).Update("count", gorm.Expr("count - ?", 1)).RowsAffected < 1 {
		tx.Rollback()
		return errors.New("no goods")
	}

	return tx.Commit().Error
}

func seckill_opt(db *gorm.DB) error {
	// 再唠叨一下，事务一旦开始，你就应该使用 tx 处理数据
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	var good Goods
	if err := tx.Debug().First(&good, "name = ?", "apple").Error; err != nil {
		tx.Rollback()
		return err
	}

	if good.Count <= 0 {
		tx.Rollback()
		return errors.New("no goods")
	}

	if tx.Model(Goods{}).Debug().Where("name = ? AND count -1  >= ?", good.Name, 0).Updates(map[string]interface{}{
		"count": gorm.Expr("count - ?", 1),
	}).RowsAffected < 0 {
		tx.Rollback()
		return errors.New("update error goods")
	}

	return tx.Commit().Error
}
