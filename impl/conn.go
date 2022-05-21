package impl

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

type Conn struct {
	C        net.Conn  // go连接对象
	IsValid  bool      // 该连接是否有效，是否关闭
	ExitChan chan bool // 退出管道通知
}

// Start 开始处理连接，接收读消息
func (c *Conn) Start() {
	fmt.Println("Conn Start...", "RemoteAddr ", c.C.RemoteAddr())
	defer c.Stop()
	// 接收连接消息
	for {
		select {
		case <-c.ExitChan:
			break
		default:
			data := make([]byte, 1024)
			_, err := c.C.Read(data)
			if err != nil {
				fmt.Println("Conn Read err", err)
				return
			}
			fmt.Println("Receive data ", string(data))

			// todo 暂时简单回显
			if err = c.SendMsg([]byte(strings.ToUpper(string(data)))); err != nil {
				fmt.Println("Conn SendMsg err ", err)
			}

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
