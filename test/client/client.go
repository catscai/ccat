package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

// DefaultHeader 实现基本的消息包
type DefaultHeader struct {
	PackType uint32
	Data     []byte
}

func (bm *DefaultHeader) Pack() ([]byte, error) {
	buff := bytes.NewBuffer([]byte{})

	if err := binary.Write(buff, binary.LittleEndian, bm.PackType); err != nil {
		return nil, err
	}

	if err := binary.Write(buff, binary.LittleEndian, bm.Data); err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

func (bm *DefaultHeader) Unpack(data []byte) error {
	reader := bytes.NewReader(data)

	if err := binary.Read(reader, binary.LittleEndian, &bm.PackType); err != nil {
		return err
	}

	//if err := binary.Read(reader, binary.LittleEndian, &bm.Data); err != nil {
	//	return err
	//}

	bm.Data = data[4:]
	return nil
}

func (bm *DefaultHeader) GetPackType() interface{} {
	return bm.PackType
}

func (bm *DefaultHeader) GetData() []byte {
	return bm.Data
}

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
		reqPack := DefaultHeader{
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
