package iface

// IWorkerGroup 工作者组
type IWorkerGroup interface {
	Init(server IServer, size uint32, queueLen uint32) // 初始化工作者组参数
	Start()                                            // 启动工作者组
	Stop()                                             // 停止，释放资源
	SendTask(request IRequest)                         // 添加处理请求任务
	GetWorkerSize() uint32                             // 获取工作者数量，工作者数量可在配置文件中配置
	SetShardWorkerFunc(f ShardWorkerHandlerFunc)       // 设置工作者选择回调,用户提供回调计算,该请求应该被分给哪个工作者执行
}

type ShardWorkerHandlerFunc func(group IWorkerGroup, request IRequest) uint32
