package iface

import "ccat/iface/imsg"

// IHeaderPackParser 解析/封装包头部
type IHeaderPackParser interface {
	HeaderPack(imsg.IHeaderPack) ([]byte, error)   // 封包，在数据包前面加上包头，封装成数据
	HeaderUnpack([]byte) (imsg.IHeaderPack, error) // 解包，解出数据包头
}
