package impl

import (
	"ccat/iface/imsg"
	"ccat/impl/msg"
	"fmt"
)

// DefaultHeaderParser 默认的header封包解包
type DefaultHeaderParser struct {
}

// HeaderUnpack 解包，解出数据包头
func (p *DefaultHeaderParser) HeaderUnpack(data []byte) (imsg.IHeaderPack, error) {
	header := msg.DefaultHeader{}
	if err := header.Unpack(data); err != nil {
		fmt.Println("HeaderUnpack Unpack err", err)
		return nil, err
	}
	return &header, nil
}
