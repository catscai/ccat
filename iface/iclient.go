package iface

import (
	"github.com/catscai/ccat/iface/imsg"
	"net"
	"time"
)

type IClient interface {
	Connection(ipVer, address string, timeOut time.Duration) error      // 连接服务器
	SetProcess(process func(conn net.Conn, msg imsg.IHeaderPack) error) // 设置消息回调
	SetDataPack(pack imsg.IDataPack)                                    // 设置处理粘包，分包
	SetHeaderOperator(parser imsg.IHeaderOperator)                      // 设置包头操作者
	Send(req, rsp imsg.IHeaderPack) error                               // 同步发送，等待请求回复;如果设置了process,那么收到消息后既会调用process回调，也会在Send接口返回
	SendASync(req imsg.IHeaderPack) error                               // 异步发送
	Valid() bool                                                        // 连接是否有效
	Close()                                                             // 关闭连接
}
