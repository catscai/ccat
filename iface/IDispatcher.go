package iface

// IDispatcher 分发消息
type IDispatcher interface {
	Dispatch(message IMessage)
}
