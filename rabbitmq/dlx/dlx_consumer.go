package main

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

func ConsumerDlx() {

	var (
		conn    *amqp.Connection
		ch      *amqp.Channel
		queue   amqp.Queue
		err     error
		delvers <-chan amqp.Delivery
		message amqp.Delivery
		ok      bool
	)

	//链接rbmq
	if conn, err = amqp.Dial("amqp://swsk33:123456@127.0.0.1:5672/"); err != nil {
		fmt.Println("amqp.Dial err: ", err)
		return
	}

	//声明信道
	if ch, err = conn.Channel(); err != nil {
		fmt.Println("conn.Channel err: ", err)
		return
	}

	//声明交换机
	if err = ch.ExchangeDeclare(
		"dlx_exchange",
		amqp.ExchangeFanout, //交换机模式fanout
		true,                //持久化
		false,               //自动删除
		false,               //是否是内置交互器,(只能通过交换器将消息路由到此交互器，不能通过客户端发送消息
		false,
		nil,
	); err != nil {
		fmt.Println("ch.ExchangeDeclare: ", err)
		return
	}

	//声明队列
	if queue, err = ch.QueueDeclare(
		"dlx_queue", //队列名称
		true,        //是否是持久化
		false,       //是否不需要确认，自动删除消息
		false,       //是否是排他队列
		false,       //是否等待服务器返回ok
		nil,
	); err != nil {
		fmt.Println("ch.QueueDeclare err: ", err)
		return
	}

	//将交换器和队列/路由key绑定
	if err = ch.QueueBind(queue.Name, "", "dlx_exchange", false, nil); err != nil {
		fmt.Println("ch.QueueBind err: ", err)
		return
	}

	//开启推模式消费
	delvers, err = ch.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		fmt.Println("ch.Consume err: ", err)
	}

	//消费接收到的消息
	for {
		select {
		case message, ok = <-delvers:
			if !ok {
				continue
			}
			go func() {
				//处理消息
				time.Sleep(time.Second * 3)
				//确认接收到的消息
				if err = message.Ack(true); err != nil {
					fmt.Println("dlx d.Ack err: ", err)
					return
				}
				fmt.Println("已确认dlx", string(message.Body))
			}()
		case <-time.After(time.Second * 1):

		}
	}
}

func main() {
	ConsumerDlx()
}
