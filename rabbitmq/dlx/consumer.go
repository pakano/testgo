package main

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

func Consumer() {
	var (
		conn            *amqp.Connection
		err             error
		ch              *amqp.Channel
		queue           amqp.Queue
		dlxExchangeName string
		delvers         <-chan amqp.Delivery
		message         amqp.Delivery
		ok              bool
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

	//设置未确认的最大消息数
	if err = ch.Qos(3, 0, false); err != nil {
		fmt.Println("ch.Qos err: ", err)
		return
	}

	dlxExchangeName = "dlx_exchange"

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

	argsQue := make(map[string]interface{})
	//添加死信队列交换器属性
	argsQue["x-dead-letter-exchange"] = dlxExchangeName
	//指定死信队列的路由key，如果为fanout,不指定使用队列路由键
	//argsQue["x-dead-letter-routing-key"] = "zhe_mess"
	//添加队列长度
	//argsQue["x-max-length"] = 1
	//添加过期时间
	argsQue["x-message-ttl"] = 60000 //单位毫秒
	//声明队列
	queue, err = ch.QueueDeclare("queue_publisher", true, false, false, false, argsQue)
	if err != nil {
		fmt.Println("ch.QueueDeclare err :", err)
		return
	}

	//绑定交换器/队列和key
	if err = ch.QueueBind(queue.Name, "key1", "exchange_publisher", false, nil); err != nil {
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
				time.Sleep(time.Second * 20)
				//确认接收到的消息
				if err = message.Ack(true); err != nil {
					//TODD: 获取到消息后，在过期时间内如果未进行确认，此消息就会流入到死信队列，此时进行消息确认就会报错
					fmt.Println("d.Ack err: ", err)
					return
				}
				fmt.Println("已确认", string(message.Body))
			}()
		case <-time.After(time.Second * 1):

		}
	}
}

func main() {
	Consumer()
}
