package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Person struct {
	Name string
	Age  int
}

var client *mongo.Client

var URI1 = "mongodb://127.0.0.1:27017"
var URI2 = "mongodb://localhost:27022,localhost:27023/?replicaSet=rs1"
var URI3 = "mongodb://localhost:27040"

func init() {
	var err error
	clientOptions := options.Client().ApplyURI(URI1)
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalln(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func insert() {
	col := client.Database("my_db").Collection("my_col")

	ret, err := col.InsertOne(context.TODO(), &Person{Name: "zhang3", Age: 10})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(ret.InsertedID)

	rets, err := col.InsertMany(context.TODO(), []interface{}{
		&Person{Name: "li4", Age: 20},
		&Person{Name: "wang5", Age: 21},
	})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(rets.InsertedIDs)
}

func FindOne() {
	fmt.Println("FindOne:")
	col := client.Database("my_db").Collection("my_col")
	var p Person

	//查询单条记录
	err := col.FindOne(context.TODO(), bson.D{{"name", "zhang3"}}).Decode(&p)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(p)
}

func Find() {
	fmt.Println("Find:")
	opt := options.Find()
	opt.SetLimit(5)

	col := client.Database("my_db").Collection("my_col")

	//Find查询多条记录，将会返回一个迭代器，使用完记得关闭
	cor, err := col.Find(context.TODO(), bson.D{{"name", "zhang3"}}, opt)
	if err != nil {
		log.Fatalln(err)
	}
	defer cor.Close(context.TODO())

	for cor.Next(context.TODO()) {
		var p Person
		cor.Decode(&p)
		fmt.Println(p)
	}
}

func updateOne() {
	//条件
	filter := bson.D{{"name", "zhang3"}}
	//如果匹配不存在，则插入
	opts := options.Update().SetUpsert(true)
	//执行的动作
	update := bson.D{
		{"$set", bson.D{
			{"age", 22}},
		},
	}

	col := client.Database("my_db").Collection("my_col")

	ret, err := col.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		log.Fatalln(err)
	}
	if ret.MatchedCount != 0 {
		fmt.Println("matched", ret.MatchedCount)
	}
	if ret.UpsertedCount != 0 {
		fmt.Println("insert", ret.UpsertedCount)
	}
}

func update() {
	//条件
	filter := bson.D{{"name", "zhang3"}}
	//如果匹配不存在，则插入
	opts := options.Update().SetUpsert(true)
	//执行的动作
	update := bson.D{
		{"$set", bson.D{
			{"age", 22}},
		},
	}

	col := client.Database("my_db").Collection("my_col")

	ret, err := col.UpdateMany(context.TODO(), filter, update, opts)
	if err != nil {
		log.Fatalln(err)
	}
	if ret.MatchedCount != 0 {
		fmt.Println("matched", ret.MatchedCount)
	}
	if ret.UpsertedCount != 0 {
		fmt.Println("insert", ret.UpsertedCount)
	}
}

func remove() {
	filter := bson.D{{"name", "zhang3"}}
	col := client.Database("my_db").Collection("my_col")
	ret, err := col.DeleteMany(context.TODO(), filter)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("remove %d records\n", ret.DeletedCount)
}

func main() {
	insert()
	FindOne()
	Find()
	updateOne()
	update()
	remove()
}
