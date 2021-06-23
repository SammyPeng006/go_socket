package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	defer fmt.Println("closed")

	server := "127.0.0.1:1024"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	fmt.Println("connect success")

	//監聽接收訊息
	go handleReader(conn)

	//主動發訊息
	var msg string
	for {
		fmt.Scanln(&msg)
		if msg == "close" || msg == "" {
			break
		}
		sender(msg, conn)
	}
}

func sender(msg string, conn net.Conn) {
	conn.Write([]byte(msg))
	//fmt.Println("send over")
}

func handleReader(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 2048)

	for {
		n, err := conn.Read(buffer)

		if err != nil {
			fmt.Println(conn.RemoteAddr().String(), " connection error: ", err)
			return
		}

		fmt.Println(conn.RemoteAddr().String(), "receive data string:\n", string(buffer[:n]))
	}
}
