package main

import (
	"fmt"

	"github.com/streadway/amqp"
)

func ProducerDlx() {
	var (
		conn *amqp.Connection
		err  error
		ch   *amqp.Channel
	)
	if conn, err = amqp.Dial("amqp://swsk33:123456@127.0.0.1:5672/"); err != nil {
		fmt.Println("amqp.Dial err :", err)
		return
	}
	defer conn.Close()

	if ch, err = conn.Channel(); err != nil {
		fmt.Println("conn.Channel err: ", err)
		return
	}

	defer ch.Close()

	//声明交换器
	if err = ch.ExchangeDeclare(
		"exchange_publisher",
		amqp.ExchangeDirect,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		fmt.Println("ch.ExchangeDeclare err: ", err)
		return
	}

	//发送消息
	if err = ch.Publish(
		"exchange_publisher",
		"key1",
		false,
		false,
		amqp.Publishing{
			Headers:      amqp.Table{},
			ContentType:  "text/plain",
			Body:         []byte("hello world dlx"),
			DeliveryMode: amqp.Persistent,
			Priority:     0,
		},
	); err != nil {
		fmt.Println("ch.Publish err: ", err)
		return
	}
}

func main() {
	ProducerDlx()
}
