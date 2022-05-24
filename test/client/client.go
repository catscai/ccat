package main

import (
	"ccat/iface/imsg"
	"ccat/impl"
	"ccat/impl/msg"
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
	client := impl.NewClient(&msg.DefaultDataPack{}, &msg.DefaultHeaderParser{}, 10, time.Millisecond*300)

	err := client.Connection("tcp4", "127.0.0.1:2233", time.Millisecond*300)
	if err != nil {
		fmt.Println("Connection failed", err)
		return
	}
	client.SetProcess(process)
	defer client.Close()
	for {
		echo := EchoMessage{
			name: "caiyanqing",
		}
		echoData, _ := echo.Pack()
		header := msg.DefaultHeader{
			PackType:  1,
			SessionID: uint64(time.Now().UnixNano()),
			Data:      echoData,
		}

		if err = client.SendASync(&header); err != nil {
			return
		}
		fmt.Println("Send data header", header, "time:", time.Now().Unix())
		time.Sleep(time.Second)
	}
}
