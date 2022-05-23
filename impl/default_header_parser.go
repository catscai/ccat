package impl

import (
	"ccat/iface/imsg"
	msg2 "ccat/impl/msg"
	"fmt"
)

type DefaultHeaderParser struct {
}

// HeaderPack 封包，在数据包前面加上包头，封装成数据
func (p *DefaultHeaderParser) HeaderPack(pack imsg.IHeaderPack) ([]byte, error) {
	return pack.Pack()
}

// HeaderUnpack 解包，解出数据包头
func (p *DefaultHeaderParser) HeaderUnpack(data []byte) (imsg.IHeaderPack, error) {
	msg := msg2.DefaultHeader{}
	if err := msg.Unpack(data); err != nil {
		fmt.Println("HeaderUnpack Unpack err", err)
		return nil, err
	}
	return &msg, nil
}
