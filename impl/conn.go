package impl

import (
	"ccat/iface"
	"errors"
	"fmt"
	"net"
)

type Conn struct {
	C        net.Conn      // go连接对象
	IsValid  bool          // 该连接是否有效，是否关闭
	ExitChan chan bool     // 退出管道通知
	Server   iface.IServer // 所属服务
}

var ConnErr = errors.New("conn error")

// Start 开始处理连接，接收读消息
func (c *Conn) Start() {
	fmt.Println("Conn Start...", "RemoteAddr ", c.C.RemoteAddr())
	defer c.Stop()
	dataPack := c.Server.GetDataPack()
	// 接收连接消息
	for {
		select {
		case <-c.ExitChan:
			break
		default:
			data, err := dataPack.ParseData(c)
			if err != nil {
				fmt.Println("Conn Start ParseData err", err)
				return
			}
			//fmt.Println("Receive data", string(data))

			// 将数据交给业务分发器分发
			c.Server.GetDispatcher().Dispatch(c, data)
		}
	}
}

// Stop 关闭连接释放资源
func (c *Conn) Stop() {
	fmt.Println("Conn Stop RemoteAddr ", c.C.RemoteAddr())
	c.IsValid = false
	c.C.Close()
	close(c.ExitChan)
}

// SendMsg 发送消息
func (c *Conn) SendMsg(data []byte) error {
	if !c.Valid() {
		return errors.New("the conn is invalid")
	}
	_, err := c.C.Write(data)

	return err
}

func (c *Conn) Valid() bool {
	return c.IsValid
}

func (c *Conn) GetConn() net.Conn {
	return c.C
}
