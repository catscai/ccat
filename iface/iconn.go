package iface

import (
	"ccat/iface/imsg"
	"net"
)

type IConn interface {
	Start()                              // 开始处理连接，接收读消息
	Stop()                               // 关闭连接释放资源
	SendMsg(pack imsg.IHeaderPack) error // 发送消息包
	SendData(data []byte) error          // 发送raw数据
	Valid() bool                         // 当前连接是否有效

	GetConn() net.Conn // 获取go连接
	GetConnID() uint32 // 获取连接ID

	SetProperty(key string, val interface{}) // 给在连接上设置属性
	GetProperty(key string) interface{}      // 获取属性
	RemoveProperty(key string)               // 删除属性
}
