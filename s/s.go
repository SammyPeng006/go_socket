package main

import (
	"container/list"
	"fmt"
	"log"
	"net"
	"os"
	uuid "satori/go.uuid"
)

var tcpConnBuff list.List

type ConnPool struct {
	uid  string
	conn net.Conn
}

func main() {
	//建立socket，監聽埠
	netListen, err := net.Listen("tcp", "localhost:1024")
	CheckError(err)
	defer netListen.Close()

	Log("Waiting for clients")
	for {
		conn, err := netListen.Accept()
		if err != nil {
			Log(conn.RemoteAddr().String(), " tcp connect error:", err)
			continue
		}
		connBuff := ConnPool{uuid.NewV4().String(), conn}
		tcpConnBuff.PushBack(connBuff)

		Log(conn.RemoteAddr().String(), " tcp connect success")
		go handleConnection(connBuff)
	}
}

//處理連線
func handleConnection(connp ConnPool) {
	defer Leave(connp)
	buffer := make([]byte, 2048)

	//通知uid
	connp.conn.Write([]byte("Your uid: " + connp.uid))

	for {
		n, err := connp.conn.Read(buffer)

		if err != nil {
			Log(connp.conn.RemoteAddr().String(), " connection error: ", err)
			return
		}

		Log(connp.conn.RemoteAddr().String(), "receive data string:\n", string(buffer[:n]))
		go SendMsg(connp, string(buffer[:n]))
	}
}

func Log(v ...interface{}) {
	log.Println(v...)
}

func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func Leave(connp ConnPool) {
	defer connp.conn.Close()

	//移除pool
	for e := tcpConnBuff.Front(); e != nil; e = e.Next() {
		if ConnPool(e.Value.(ConnPool)).uid == connp.uid {
			tcpConnBuff.Remove(e)
			break
		}
	}

	//通知離開
	leaveUid := connp.uid
	for e := tcpConnBuff.Front(); e != nil; e = e.Next() {
		ConnPool(e.Value.(ConnPool)).conn.Write([]byte(leaveUid + " already leaved"))
	}
}

func SendMsg(connp ConnPool, msg string) {

	//通知
	sayerUid := connp.uid
	for e := tcpConnBuff.Front(); e != nil; e = e.Next() {
		cp := ConnPool(e.Value.(ConnPool))
		if cp.uid != sayerUid {
			cp.conn.Write([]byte(sayerUid + " : " + msg))
		}
	}
}
