package main

import (
	"ccat"
	"ccat/iface"
	"ccat/impl"
	"fmt"
)

type EchoMessage struct {
	impl.BaseMessage
	name string
}

func (em *EchoMessage) Unpack(data []byte) error {
	em.name = string(data)
	return nil
}

func (em *EchoMessage) Pack() ([]byte, error) {
	return []byte(em.name), nil
}

func Deal(conn iface.IConn, message iface.IMessage) error {
	fmt.Println("deal recv message start")
	defer fmt.Println("deal recv message end")
	req := message.(*EchoMessage)
	fmt.Printf("Recv message:%+v\n", *req)
	return nil
}

func main() {
	service := ccat.GetServer("tcp_test")
	if service == nil {
		fmt.Println("get service is nil")
		return
	}
	service.GetDispatcher().RegisterHandler(uint32(1), &EchoMessage{}, Deal)
	ccat.Run()
	// todo 接下来 开发工作任务池
}
