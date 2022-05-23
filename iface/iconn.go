package iface

import "net"

type IConn interface {
	Start()               // 开始处理连接，接收读消息
	Stop()                // 关闭连接释放资源
	SendMsg([]byte) error // 发送消息
	Valid() bool          // 当前连接是否有效

	GetConn() net.Conn // 获取go连接
	GetConnID() uint32 // 获取连接ID
}
