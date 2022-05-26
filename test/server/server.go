package main

import (
	"ccat"
	"ccat/iface"
	"ccat/iface/imsg"
	"ccat/impl/msg"
	"ccat/test"
	"fmt"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
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

func Deal(ctx *iface.CatContext, request iface.IRequest, message imsg.IMessage) error {
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

func DealSimpleMessage(ctx *iface.CatContext, reqMsg, rspMsg imsg.IMessage) error {
	req := reqMsg.(*EchoMessage)
	rsp := rspMsg.(*EchoMessage)
	fmt.Println("DealSimpleMessage recv req", *req)

	rsp.name = req.name + "HHHHH---KKK"

	return nil
}

func DealSimplePB(ctx *iface.CatContext, reqMsg, rspMsg proto.Message) error {
	req := reqMsg.(*test.TestRQ)
	rsp := rspMsg.(*test.TestRS)

	fmt.Println("DealSimplePB recv req", *req)

	rsp.Uid = req.Uid
	rsp.Name = req.Name
	rsp.Age = req.Age
	rsp.Reply = proto.String("i love you")
	ctx.Info("DealSimplePB recv data", zap.Any("req", *req), zap.Any("rsp", *rsp))
	return nil
}

const (
	TestPackType1RQ uint32 = 1
	TestPackType1RS uint32 = 2

	TestPackType2RQ uint32 = 3
	TestPackType2RS uint32 = 4

	TestPackType3RQ uint32 = 5
	TestPackType4RS uint32 = 6
)

func main() {
	service := ccat.GetServer("tcp_test")
	if service == nil {
		fmt.Println("get service is nil")
		return
	}

	service.GetDispatcher().RegisterHandler(TestPackType1RQ, &EchoMessage{}, Deal)
	service.GetDispatcher().RegisterHandlerSimple(TestPackType2RQ, TestPackType2RS, &EchoMessage{}, &EchoMessage{}, DealSimpleMessage)
	service.GetDispatcher().RegisterHandlerSimplePB(TestPackType3RQ, TestPackType4RS, &test.TestRQ{}, &test.TestRS{}, DealSimplePB)
	ccat.Run()
}
