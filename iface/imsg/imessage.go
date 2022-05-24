package imsg

// IMessage 消息接口， 提供封包和解包操作;
type IMessage interface {
	Pack() ([]byte, error) // 消息封包
	Unpack([]byte) error   // 消息解包
}
