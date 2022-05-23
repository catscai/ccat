package iface

// IRequest 表示一个请求， 包含连接信息，和数据包
type IRequest interface {
	GetConn() IConn             // 获取连接信息
	GetHeaderPack() IHeaderPack // 获取头数据包
}
