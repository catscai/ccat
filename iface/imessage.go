package iface

// IMessage 消息接口， 提供封包和解包操作;
type IMessage interface {
	Pack() ([]byte, error)
	Unpack([]byte) error
}
