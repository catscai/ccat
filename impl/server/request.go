package server

import (
	"ccat/iface"
	"ccat/iface/imsg"
)

type Request struct {
	Conn       iface.IConn
	HeaderPack imsg.IHeaderPack
}

// GetConn 获取连接信息
func (r *Request) GetConn() iface.IConn {
	return r.Conn
}

// GetHeaderPack 获取头数据包
func (r *Request) GetHeaderPack() imsg.IHeaderPack {
	return r.HeaderPack
}
