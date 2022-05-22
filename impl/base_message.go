package impl

// BaseMessage 继承消息接口，没有具体实现
type BaseMessage struct {
}

func (b *BaseMessage) Pack() ([]byte, error) {
	return nil, nil
}

func (b *BaseMessage) Unpack([]byte) error {
	return nil
}
