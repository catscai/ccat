package iface

// IServer 服务器接口
type IServer interface {
	Start()
	Stop()
	Run()

	GetDataPack() IDataPack         // 获取数据包处理对象
	SetDataPack(packDeal IDataPack) // 设置数据包处理对象
}
