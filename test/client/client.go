package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp4", "127.0.0.1:2233")
	if err != nil {
		fmt.Println("connection err", err)
		return
	}
	for {
		_, err = conn.Write([]byte("hello world"))
		if err != nil {
			fmt.Println("Send data err", err)
			break
		}
		data := make([]byte, 1024)
		_, err = conn.Read(data)
		if err != nil {
			fmt.Println("Recv data err", err)
			break
		}
		fmt.Println("Recv data", string(data))
		time.Sleep(time.Second)
	}
}
