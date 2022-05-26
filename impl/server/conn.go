package server

import (
	"ccat/clog"
	"ccat/config"
	"ccat/iface"
	"ccat/iface/imsg"
	"context"
	"errors"
	"go.uber.org/zap"
	"net"
	"sync"
)

type Conn struct {
	clog.ICatLog
	C       net.Conn        // go连接对象
	ConnID  uint32          // 连接id
	IsValid bool            // 该连接是否有效，是否关闭
	ctx     context.Context // 用来处理关闭
	cancel  context.CancelFunc
	Server  iface.IServer // 所属服务

	propertyMap   map[string]interface{} // 属性map
	propertyMutex sync.RWMutex           // 属性操作锁
}

var ConnErr = errors.New("conn error")

func NewConnection(c net.Conn, connID uint32, ser iface.IServer) *Conn {
	catConn := &Conn{
		ICatLog:     clog.AppLogger().Clone(),
		C:           c,
		IsValid:     true,
		Server:      ser,
		ConnID:      connID,
		propertyMap: nil,
	}

	return catConn
}

func (c *Conn) StartReader() {
	c.Info("Conn Start...", zap.Any("RemoteAddr ", c.C.RemoteAddr()))
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
				c.Error("Conn Start ParseData failed", zap.Any("err", err))
				return
			}
			cfg := config.GetBaseServiceCfg(c.Server.GetName())
			if cfg.MaxPackLen > 0 && uint32(len(data)) > cfg.MaxPackLen {
				c.Warn("ParseData Recv packlen over max pack length limit", zap.Int("packLen", len(data)))
				return
			}
			//fmt.Println("Receive data", string(data))
			// 解析出包头
			//header, err := c.Server.GetHeaderParser().HeaderUnpack(data)
			header := c.Server.GetHeaderOperator().Get()
			err = header.Unpack(data)
			if err != nil {
				c.Error("Conn HeaderUnpack failed", zap.Any("err", err))
				return
			}
			r := Request{
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

	// 调用连接开始时的回调
	c.Server.CallConnectStart(c)

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
		c.Error("Conn SendMsg pack.Pack failed", zap.Any("err", err))
		return err
	}
	data, err := c.Server.GetDataPack().ReorganizeData(packData)
	if err != nil {
		c.Error("Conn SendMsg ReorganizeData failed", zap.Any("err", err))
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
	c.Info("Conn release RemoteAddr ", zap.Any("remoteAddr", c.C.RemoteAddr()))
	c.IsValid = false
	c.C.Close()
	// 从连接管理器中删除
	c.Server.GetConnManager().Remove(c)

	// 调用连接退出时回调
	c.Server.CallConnectEnd(c)
}

// SetProperty 给在连接上设置属性
func (c *Conn) SetProperty(key string, val interface{}) {
	c.propertyMutex.Lock()
	defer c.propertyMutex.Unlock()
	if c.propertyMap == nil {
		c.propertyMap = make(map[string]interface{})
	}
	c.propertyMap[key] = val
}

// GetProperty 获取属性
func (c *Conn) GetProperty(key string) interface{} {
	c.propertyMutex.RLock()
	defer c.propertyMutex.RUnlock()
	if c.propertyMap == nil {
		return nil
	}
	if val, ok := c.propertyMap[key]; ok {
		return val
	}
	return nil
}

// RemoveProperty 删除属性
func (c *Conn) RemoveProperty(key string) {
	c.propertyMutex.Lock()
	defer c.propertyMutex.Unlock()
	if c.propertyMap != nil {
		delete(c.propertyMap, key)
	}
}
