package iface

import (
	"ccat/iface/imsg"
)

// IDispatcher 分发消息
type IDispatcher interface {
	Dispatch(request IRequest)                                                        // 将消息分发给处理函数
	RegisterHandler(packType interface{}, message imsg.IMessage, deal MsgHandlerFunc) // 注册消息回调
	Remove(packType interface{})                                                      // 删除回调映射关系
}

type MsgHandlerFunc func(conn IConn, iMessage imsg.IMessage) error
