package iface

import "net"

// IDataPack 负责收到数据时解包，和发送数据时封包；这里处理的是流数据,不是消息包
type IDataPack interface {
	ParseData(conn net.Conn) ([]byte, error)    // 解析读数据,返回包数据
	ReorganizeData(data []byte) ([]byte, error) // 重新组织数据，加上头部包长度
}
