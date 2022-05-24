package imsg

// IHeaderOperator 当用户不想使用默认的包头时(DefaultHeader), 可以实现这两个接口
type IHeaderOperator interface {
	Get() IHeaderPack                                                     // 获取一个包头对象
	Full(msgType interface{}, data []byte, other IHeaderPack) IHeaderPack // 填充一个包头对象
}
