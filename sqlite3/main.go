package main

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Person struct {
	ID   int `gorm:"primaryKey"`
	Name string
}

func main() {
	db, err := gorm.Open(sqlite.Open("./db"), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	db.AutoMigrate(Person{})
	var ps []Person
	db.Debug().Create(&ps)
}
