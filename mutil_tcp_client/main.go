package main

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"net"
	"sync"
	"test/util"
	"time"
)

type RQ struct {
	Session uint32
}

type RP struct {
	Session uint32
}

type Msg struct {
	data *RQ
	ch   chan *RP
}

type TCPClient struct {
	deviceid string
	conn     net.Conn
	mu       sync.Mutex
	tokens   chan chan *RP
	msg      chan Msg
	rch      chan *RP
	sch      chan *RQ
	lastSend time.Time
}

func NewTCPClient(deviceid string, conn net.Conn, maxclientcount int) *TCPClient {
	c := &TCPClient{deviceid: deviceid, lastSend: time.Now()}
	c.conn = conn
	c.tokens = make(chan chan *RP, maxclientcount)
	c.msg = make(chan Msg, maxclientcount)
	for i := 0; i < maxclientcount; i++ {
		c.tokens <- make(chan *RP)
	}
	c.rch = make(chan *RP) //接收通道
	c.sch = make(chan *RQ) //发送通道

	return c
}

func (client *TCPClient) recvRP() (*RP, error) {
	var buf = make([]byte, 4)
	_, err := io.ReadFull(client.conn, buf)
	if err != nil {
		return nil, err
	}
	n := binary.BigEndian.Uint32(buf)
	data := make([]byte, n)
	_, err = io.ReadFull(client.conn, data)
	if err != nil {
		return nil, err
	}
	var rp RP
	err = json.Unmarshal(data, &rp)
	if err != nil {
		return nil, err
	}
	return &rp, nil
}

func (client *TCPClient) sendRQ(rq *RQ) error {
	client.mu.Lock()
	defer client.mu.Unlock()
	data, err := json.Marshal(rq)
	if err != nil {
		return err
	}

	//发送指令最大超时
	client.conn.SetWriteDeadline(time.Now().Add(time.Second * 5))

	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(len(data)))
	_, err = client.conn.Write(buf)
	if err != nil {
		return err
	}
	_, err = client.conn.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (client *TCPClient) run() {
	defer util.MainDefer()
	go func() {
		defer util.MainDefer()
		defer log.Println("debug:", "sch exit........")
		for {
			rq := <-client.sch
			if rq == nil {
				return
			}

			err := client.sendRQ(rq)
			if err != nil {
				log.Println("err:", err)
				return
			}
		}
	}()
	go func() {
		defer util.MainDefer()
		defer log.Println("debug:", "rch exit........")
		for {
			rp, err := client.recvRP()
			if err != nil {
				if err != io.EOF {
					log.Println("err:", err)
				}
				return
			}
			client.rch <- rp
		}
	}()
	client.manage(client.sch, client.rch)
}

func (client *TCPClient) manage(sch chan<- *RQ, rch <-chan *RP) {
	players := make(map[uint32]chan *RP)
	var id uint32
	for {
		select {
		case msg := <-client.msg:
			id++
			msg.data.Session = id
			players[id] = msg.ch
			sch <- msg.data
		case data := <-rch:
			if data == nil {
				log.Println("debug:", "manage exit........")
				return
			}
			session := data.Session
			ch, ok := players[session]
			if ok {
				ch <- data
				delete(players, session)
			} else {
				log.Println("err:", "not found session ch")
			}
		}
	}
}

func (client *TCPClient) Exit() {
	log.Println("debug:", "device:", client.deviceid, "exit")
	client.conn.Close()
	//TCPManager.clients.Delete(client.deviceid)
	client.rch <- nil
	client.sch <- nil
}

// 返回proto中的数据结构
func (client *TCPClient) SendAndReceive(rq *RQ) *RP {
	//客户端发送队列满
	if len(client.tokens) == 0 {
		return nil
	}

	ch := <-client.tokens
	// 回收信道
	defer func() { client.tokens <- ch }()

	client.msg <- Msg{rq, ch}

	//最大超时5s
	select {
	case rp := <-ch:
		return rp
	case <-time.After(time.Second * 5):
		log.Println("err:", client.deviceid, "timeout")
		return nil
	}
}
