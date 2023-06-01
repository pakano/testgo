package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("./main localaddr remoteaddr")
	}
	localaddr := os.Args[1]
	remoteaddr := os.Args[2]

	l, err := net.Listen("tcp", localaddr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("err:", err)
		} else {
			go Serve(conn, remoteaddr)
		}
	}
}

func Serve(conn1 net.Conn, remoteAddr string) {
	defer conn1.Close()

	conn2, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		log.Println("err:", err)
	}
	defer conn2.Close()

	go io.Copy(conn1, conn2)
	io.Copy(conn2, conn1)
}
