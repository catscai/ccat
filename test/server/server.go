package main

import (
	"ccat"
	"ccat/iface"
	"ccat/iface/imsg"
	"ccat/impl/msg"
	"fmt"
)

type EchoMessage struct {
	msg.BaseMessage
	name string
}

func (em *EchoMessage) Unpack(data []byte) error {
	em.name = string(data)
	return nil
}

func (em *EchoMessage) Pack() ([]byte, error) {
	return []byte(em.name), nil
}

func Deal(request iface.IRequest, message imsg.IMessage) error {
	fmt.Println("deal recv message start")
	defer fmt.Println("deal recv message end")
	req := message.(*EchoMessage)
	fmt.Printf("Recv message:%+v\n", *req)

	// 回复
	echo := EchoMessage{
		name: "yunshuipiao",
	}
	echoData, _ := echo.Pack()
	header := msg.DefaultHeader{
		PackType:  2,
		SessionID: request.GetHeaderPack().GetSessionID().(uint64),
		Data:      echoData,
	}
	request.GetConn().SendMsg(&header)
	fmt.Println("deal SendMsg", header)
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
