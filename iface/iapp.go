package iface

// IApp 整个项目
type IApp interface {
	AddServer(server IServer)      // 向整个项目添加server,如tcp,udp,rpc
	GetServer(name string) IServer // 根据服务名字获取server对象
	Run()
}
