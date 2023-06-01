// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	u := url.URL{Scheme: "ws", Host: ":8080", Path: "/echo"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	c.SetPingHandler(func(appData string) error {
		fmt.Println("recv ping")
		c.WriteMessage(websocket.PongMessage, nil)
		return nil
	})

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("recv: %s", message)
	}
}
