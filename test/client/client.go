package main

import (
	"ccat/iface/imsg"
	"ccat/impl"
	"ccat/impl/msg"
	"ccat/test"
	"fmt"
	"github.com/golang/protobuf/proto"
	"net"
	"time"
)

type EchoMessage struct {
	name string
}

func (em *EchoMessage) Unpack(data []byte) error {
	em.name = string(data)
	return nil
}

func (em *EchoMessage) Pack() ([]byte, error) {
	return []byte(em.name), nil
}

func process(conn net.Conn, pack imsg.IHeaderPack) error {
	header := pack.(*msg.DefaultHeader)

	fmt.Println("process recv data header", header)

	echo := EchoMessage{}
	if err := echo.Unpack(header.GetData()); err != nil {
		fmt.Println("process echo.Unpack err", err)
		return err
	}
	fmt.Println("process recv echo", echo)
	return nil
}

func main() {
	client := impl.NewClient(&msg.DefaultDataPack{}, &msg.DefaultHeaderOperator{}, 10, time.Millisecond*300)

	err := client.Connection("tcp4", "127.0.0.1:2233", time.Second)
	if err != nil {
		fmt.Println("Connection failed", err)
		return
	}
	client.SetProcess(process)
	defer client.Close()
	for {
		req := test.TestRQ{
			Uid:  proto.Uint64(100010),
			Name: proto.String("caiyanqing"),
			Age:  proto.Uint32(25),
		}
		reqData, err := proto.Marshal(&req)
		if err != nil {
			fmt.Println("proto.Marshal err", err)
			return
		}
		header := msg.DefaultHeader{
			PackType:  1,
			SessionID: uint64(time.Now().UnixNano()),
			Data:      reqData,
		}

		//if err = client.SendASync(&header); err != nil {
		//	return
		//}
		rsp := msg.DefaultHeader{}
		if err = client.Send(&header, &rsp); err != nil {
			fmt.Println("Send err", err)
			return
		}
		rs := test.TestRS{}
		if err = proto.Unmarshal(rsp.GetData(), &rs); err != nil {
			fmt.Println("proto.Unmarshal err", err)
			return
		}
		fmt.Printf("Send recv uid:%d,name:%s,age:%d,reply:%s\n",
			rs.GetUid(), rs.GetName(), rs.GetAge(), rs.GetReply())
		fmt.Println("Send data header", header, "time:", time.Now().Unix())
		time.Sleep(time.Second)
	}
}
