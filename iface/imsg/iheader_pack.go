package imsg

// IHeaderPack 包头
type IHeaderPack interface {
	Pack() ([]byte, error)    // 消息封包
	Unpack([]byte) error      // 消息解包
	GetPackType() interface{} // 获取包类型
	GetData() []byte          // 获取包中数据
}
