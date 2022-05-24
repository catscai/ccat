package iface

import (
	"ccat/iface/imsg"
)

// IServer 服务器接口
type IServer interface {
	Start() // 创建socket启动监听
	Stop()  // 释放资源
	Run()   // 调用Start 同时执行一些其他处理代码

	SetDataPack(packDeal imsg.IDataPack) // 设置数据包处理对象
	GetDataPack() imsg.IDataPack         // 获取数据包处理对象

	SetDispatcher(dispatcher IDispatcher) // 设置消息分发器
	GetDispatcher() IDispatcher           // 获取消息分发器

	SetHeaderParser(parser imsg.IHeaderPackParser) // 设置包头解析器
	GetHeaderParser() imsg.IHeaderPackParser       // 获取包头解析器

	GetWorkerGroup() IWorkerGroup // 获取工作者组

	GetConnManager() IConnManager // 获取连接管理器

	SetConnectStartHandler(ConnStatusHandler) // 设置连接建立成功后的回调
	SetConnectEndHandler(ConnStatusHandler)   // 设置连接关闭时的回调
	CallConnectStart(conn IConn)              // 调用连接开始回调
	CallConnectEnd(conn IConn)                // 调用连接退出回调

	GetAddr() string // 获取服务监听地址
	GetName() string // 获取服务名
}

type ConnStatusHandler func(server IServer, conn IConn) error
