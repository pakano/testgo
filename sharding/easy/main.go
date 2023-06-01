package main

import (
	"fmt"
	"hash/crc32"
	"os"
	"test/util"

	"github.com/bwmarrin/snowflake"
	"gorm.io/gorm"
)

const (
	MAXTABLENUM = 5
)

type Order struct {
	ID   int64  `gorm:"type:unsigned bigint,primaryKey"`
	Name string `gorm:"type:varchar(20)"`
	Data string `gorm:"type:varchar(100)"`
}

type OrderIdx struct {
	Name  string `gorm:"type:varchar(20)"`
	OldID int64  `gorm:"type:unsigned bigint"`
}

var node *snowflake.Node

func init() {
	var err error
	node, err = snowflake.NewNode(1)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func main() {
	db, err := util.NewMysqlInstance("root:@tcp(127.0.0.1:4000)/test?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < MAXTABLENUM; i++ {
		table := fmt.Sprintf("order_%02d", i)
		db.Exec(`DROP TABLE IF EXISTS ` + table)
		db.Exec(`CREATE TABLE ` + table + ` (
			id bigint unsigned primary key,
			name varchar(20),
			data varchar(100)
		)`)
	}

	for i := 0; i < MAXTABLENUM; i++ {
		table := fmt.Sprintf("order_idx_%02d", i)
		db.Exec(`DROP TABLE IF EXISTS ` + table)
		db.Exec(`CREATE TABLE ` + table + ` (
			name varchar(20),
			old_id  bigint unsigned,
			Key idx_name(name)
		)`)
	}

	for {
		id := NewOrder(db)
		order1 := GetOrder(db, id)
		order2 := GetOrderByName(db, order1.Name)
		fmt.Println(order1, order2)
		DelOrder(db, order1.ID)
	}
}

func NewOrder(db *gorm.DB) int64 {
	ID := node.Generate()
	Name := util.GetRandomToken(5)

	tableNum := ID.Int64() % MAXTABLENUM
	table := fmt.Sprintf("order_%02d", tableNum)

	db.Debug().Table(table).Create(Order{ID: ID.Int64(), Name: Name})

	newOrderIdx(db, Order{ID: ID.Int64(), Name: Name})

	return int64(ID)
}

func newOrderIdx(db *gorm.DB, order Order) {
	tableNum := crc32.ChecksumIEEE([]byte(order.Name)) % MAXTABLENUM
	table := fmt.Sprintf("order_idx_%02d", tableNum)
	db.Debug().Table(table).Create(OrderIdx{Name: order.Name, OldID: order.ID})
}

func delOrderIdx(db *gorm.DB, order Order) {
	tableNum := crc32.ChecksumIEEE([]byte(order.Name)) % MAXTABLENUM
	table := fmt.Sprintf("order_idx_%02d", tableNum)
	db.Debug().Table(table).Delete(nil, "name = ?", order.Name)
}

func GetOrderIdx(db *gorm.DB, name string) int64 {
	tableNum := crc32.ChecksumIEEE([]byte(name)) % MAXTABLENUM
	table := fmt.Sprintf("order_idx_%02d", tableNum)

	var orderIdx OrderIdx
	db.Debug().Table(table).First(&orderIdx, "name = ?", name)
	return orderIdx.OldID
}

func GetOrder(db *gorm.DB, id int64) Order {
	tableNum := id % MAXTABLENUM
	table := fmt.Sprintf("order_%02d", tableNum)

	var order Order
	db.Debug().Table(table).First(&order, "id = ?", id)
	return order
}

func DelOrder(db *gorm.DB, id int64) {
	order := GetOrder(db, id)

	tableNum := id % MAXTABLENUM
	table := fmt.Sprintf("order_%02d", tableNum)

	db.Debug().Table(table).Delete(nil, "id = ?", id)

	delOrderIdx(db, order)
}

func GetOrderByName(db *gorm.DB, name string) Order {
	ID := GetOrderIdx(db, name)
	order := GetOrder(db, ID)
	return order
}
