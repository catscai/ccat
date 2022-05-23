package iface

// IServer 服务器接口
type IServer interface {
	Start() // 创建socket启动监听
	Stop()  // 释放资源
	Run()   // 调用Start 同时执行一些其他处理代码

	SetDataPack(packDeal IDataPack) // 设置数据包处理对象
	GetDataPack() IDataPack         // 获取数据包处理对象

	SetDispatcher(dispatcher IDispatcher) // 设置消息分发器
	GetDispatcher() IDispatcher           // 获取消息分发器

	SetHeaderParser(parser IHeaderPackParser) // 设置包头解析器
	GetHeaderParser() IHeaderPackParser       // 获取包头解析器

	GetWorkerGroup() IWorkerGroup // 获取工作者组

	GetAddr() string // 获取服务监听地址
	GetName() string // 获取服务名
}
