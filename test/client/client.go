package main

import (
	"ccat/impl/msg"
	"encoding/binary"
	"fmt"
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

func main() {
	conn, err := net.Dial("tcp4", "127.0.0.1:2233")
	if err != nil {
		fmt.Println("connection err", err)
		return
	}
	defer conn.Close()
	for {
		req := EchoMessage{
			name: "hello caiyanqing",
		}
		reqData, err := req.Pack()
		if err != nil {
			fmt.Println("req pack err", err)
			break
		}
		reqPack := msg.DefaultHeader{
			PackType: 1,
			Data:     reqData,
		}
		reqPackData, err := reqPack.Pack()
		if err != nil {
			fmt.Println("send pack err", err)
			break
		}
		packLen := len(reqPackData)
		data := make([]byte, 4)
		binary.LittleEndian.PutUint32(data, uint32(packLen))
		data = append(data, reqPackData...)
		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("Send data err", err)
			break
		}
		fmt.Printf("sending data:%+v\n", req)
		//data := make([]byte, 1024)
		//_, err = conn.Read(data)
		//if err != nil {
		//	fmt.Println("Recv data err", err)
		//	break
		//}
		//fmt.Println("Recv data", string(data))
		time.Sleep(time.Second)
	}
}
