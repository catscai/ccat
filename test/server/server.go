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

func DealSimpleMessage(reqMsg, rspMsg imsg.IMessage) error {
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

func main() {
	service := ccat.GetServer("tcp_test")
	if service == nil {
		fmt.Println("get service is nil")
		return
	}

	//service.GetDispatcher().RegisterHandler(uint32(1), &EchoMessage{}, Deal)
	//service.GetDispatcher().RegisterHandlerSimple(uint32(1), uint32(2), &EchoMessage{}, &EchoMessage{}, DealSimpleMessage)
	service.GetDispatcher().RegisterHandlerSimplePB(uint32(1), uint32(2), &test.TestRQ{}, &test.TestRS{}, DealSimplePB)
	ccat.Run()
	// todo 接下来 开发工作任务池
}
