package iface

// IApp 整个项目
type IApp interface {
	AddServer(server IServer) // 向整个项目添加server,如tcp,udp,rpc
	Run()
}
