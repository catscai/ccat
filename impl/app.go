package impl

import (
	"ccat/iface"
)

type App struct {
	Servers []iface.IServer
}

// AddServer 向整个项目添加server,如tcp,udp,rpc
func (app *App) AddServer(server iface.IServer) {
	app.Servers = append(app.Servers, server)
}

func (app *App) Run() {
	for i := 0; i < len(app.Servers); i++ {
		app.Servers[i].Run()
	}
	// 阻塞这里，不让进程退出
	select {}
}
