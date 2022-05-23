package iface

// IConnManager 连接管理器
type IConnManager interface {
	Add(conn IConn)
	Remove(conn IConn)
	RemoveAndClose(connID uint32)
	Get(connID uint32) IConn
	GetSize() uint32
	Clear()
}
