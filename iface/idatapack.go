package iface

// IDataPack 负责收到数据时解包，和发送数据时封包；这里处理的是流数据,不是消息包
type IDataPack interface {
	ParseData(conn IConn) ([]byte, error)        // 解析读数据,返回包数据
	ReorganizeData(msg IMessage) ([]byte, error) // 重新组织数据，将一个消息组织成要发送的data
}
