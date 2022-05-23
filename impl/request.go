package impl

import "ccat/iface"

type Request struct {
	Conn       iface.IConn
	HeaderPack iface.IHeaderPack
}

// GetConn 获取连接信息
func (r *Request) GetConn() iface.IConn {
	return r.Conn
}

// GetHeaderPack 获取头数据包
func (r *Request) GetHeaderPack() iface.IHeaderPack {
	return r.HeaderPack
}
