package server

import (
	"ccat/config"
	"ccat/iface"
	"ccat/iface/imsg"
	"ccat/impl"
	"context"
	"errors"
	"fmt"
	"net"
)

type Conn struct {
	C       net.Conn        // go连接对象
	ConnID  uint32          // 连接id
	IsValid bool            // 该连接是否有效，是否关闭
	ctx     context.Context // 用来处理关闭
	cancel  context.CancelFunc
	Server  iface.IServer // 所属服务
}

var ConnErr = errors.New("conn error")

func NewConnection(c net.Conn, connID uint32, ser iface.IServer) *Conn {
	catConn := &Conn{
		C:       c,
		IsValid: true,
		Server:  ser,
		ConnID:  connID,
	}

	return catConn
}

func (c *Conn) StartReader() {
	fmt.Println("Conn Start...", "RemoteAddr ", c.C.RemoteAddr())
	defer c.Stop()
	// 接收连接消息
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			// 处理tcp粘包
			data, err := c.Server.GetDataPack().ParseData(c.C)
			if err != nil {
				fmt.Println("Conn Start ParseData err", err)
				return
			}
			cfg := config.GetBaseServiceCfg(c.Server.GetName())
			if cfg.MaxPackLen > 0 && uint32(len(data)) > cfg.MaxPackLen {
				fmt.Println("ParseData Recv packlen over max pack length limit, packLen", uint32(len(data)))
				return
			}
			//fmt.Println("Receive data", string(data))
			// 解析出包头
			header, err := c.Server.GetHeaderParser().HeaderUnpack(data)
			if err != nil {
				fmt.Println("Conn HeaderUnpack err", err)
				return
			}
			r := impl.Request{
				Conn:       c,
				HeaderPack: header,
			}
			// 将数据交给业务处理工作者组
			c.Server.GetWorkerGroup().SendTask(&r)
		}
	}
}

// Start 开始处理连接，接收读消息
func (c *Conn) Start() {
	c.ctx, c.cancel = context.WithCancel(context.Background())

	go c.StartReader()

	select {
	case <-c.ctx.Done():
		c.release()
		return
	}
}

// Stop 关闭连接释放资源
func (c *Conn) Stop() {
	c.cancel() // 发送取消信号
}

// SendMsg 发送消息包
func (c *Conn) SendMsg(pack imsg.IHeaderPack) error {
	packData, err := pack.Pack()
	if err != nil {
		fmt.Println("Conn SendMsg pack.Pack err", err)
		return err
	}
	data, err := c.Server.GetDataPack().ReorganizeData(packData)
	if err != nil {
		fmt.Println("Conn SendMsg ReorganizeData err", err)
		return err
	}
	return c.SendData(data)
}

// SendData 发送raw数据
func (c *Conn) SendData(data []byte) error {
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

// GetConnID 获取连接ID
func (c *Conn) GetConnID() uint32 {
	return c.ConnID
}

// 释放资源
func (c *Conn) release() {
	fmt.Println("Conn release RemoteAddr ", c.C.RemoteAddr())
	c.IsValid = false
	c.C.Close()
	// 从连接管理器中删除
	c.Server.GetConnManager().Remove(c)
}
