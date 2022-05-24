package msg

import (
	"ccat/iface/imsg"
	"fmt"
)

// DefaultHeaderParser 默认的header封包解包
type DefaultHeaderParser struct {
}

// HeaderUnpack 解包，解出数据包头
func (p *DefaultHeaderParser) HeaderUnpack(data []byte) (imsg.IHeaderPack, error) {
	header := DefaultHeader{}
	if err := header.Unpack(data); err != nil {
		fmt.Println("HeaderUnpack Unpack err", err)
		return nil, err
	}
	return &header, nil
}
