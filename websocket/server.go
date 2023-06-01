package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 5 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 10 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{} // use default options

func handler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	fmt.Println(1)

	go reader(c)
	go writer(c)
}

func reader(conn *websocket.Conn) {
	defer conn.Close()

	conn.SetPongHandler(func(appData string) error {
		fmt.Println("recv pong")
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	conn.SetReadDeadline(time.Now().Add(pongWait))

	for {
		mt, d, err := conn.ReadMessage()
		if err != nil {
			log.Println("err:", err)
			break
		}
		fmt.Println(mt, string(d))
		conn.WriteMessage(websocket.TextMessage, []byte("hello"))
	}
}

func writer(conn *websocket.Conn) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		conn.Close()
	}()

	for {
		select {
		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println("err1:", err)
				return
			}
		}
	}
}

func main() {
	http.HandleFunc("/echo", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// //中间件
// redis mysql es clickhouse kalfka
// //运维
// nginx docker k8s
// //go 技术栈
// gin echo kit corba
