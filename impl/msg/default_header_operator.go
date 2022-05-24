package msg

import "ccat/iface/imsg"

// DefaultHeaderOperator 默认的包头操作者,如果用户不想使用默认的包头来解析,只需要自定义包头,然后重写操作者设置到server中即可
type DefaultHeaderOperator struct {
}

// Get 获取一个包头对象
func (o *DefaultHeaderOperator) Get() imsg.IHeaderPack {
	return &DefaultHeader{}
}

// Full 填充一个包头对象
func (o *DefaultHeaderOperator) Full(msgType interface{}, bodyData []byte, other imsg.IHeaderPack) imsg.IHeaderPack {
	return &DefaultHeader{
		PackType:  msgType.(uint32),
		SessionID: other.GetSessionID().(uint64),
		Data:      bodyData,
	}
}
