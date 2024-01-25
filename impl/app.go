package impl

import (
	"github.com/catscai/ccat/clog"
	"github.com/catscai/ccat/iface"
)

type App struct {
	clog.ICatLog
	ServerMap map[string]iface.IServer
}

// AddServer 向整个项目添加server,如tcp,udp,rpc
func (app *App) AddServer(server iface.IServer) {
	app.ServerMap[server.GetName()] = server
}

// GetServer 根据服务名字获取server对象
func (app *App) GetServer(name string) iface.IServer {
	if ser, ok := app.ServerMap[name]; ok {
		return ser
	}

	return nil
}

func (app *App) Run() {
	for name := range app.ServerMap {
		app.ServerMap[name].Run()
	}
	// 阻塞这里，不让进程退出
	select {}
}
