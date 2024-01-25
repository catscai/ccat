package server

import (
	"fmt"
	"github.com/catscai/ccat/clog"
	"github.com/catscai/ccat/config"
	"github.com/catscai/ccat/iface"
	"github.com/catscai/ccat/impl"
	"github.com/catscai/ccat/impl/msg"
	server2 "github.com/catscai/ccat/impl/server"
)

// 全局项目app
var gApp *impl.App

func init() {
	// 加载全局配置
	if err := config.Reload(); err != nil {
		panic(fmt.Sprintf("load config err:%+v", err))
	}

	gApp = &impl.App{
		ServerMap: make(map[string]iface.IServer),
	}
	clog.InitAppLogger() // 初始化全局日志
	if config.AppCfg.IsTcpService {
		tcpCfg := &config.AppCfg.TcpCfg
		for name, info := range *tcpCfg {
			if info.Auto { // 自动创建服务监听
				AddServer(NewTcpService(name, "tcp4", info.IP, info.Port))
			}
		}
	}
}

func AddServer(server iface.IServer) {
	gApp.AddServer(server)
}

func GetServer(name string) iface.IServer {
	return gApp.GetServer(name)
}

func Run() {
	gApp.Run()
}

func NewTcpService(name, ipVer, ip string, port uint32) iface.IServer {
	ser := &server2.TcpService{
		ICatLog:     clog.AppLogger().Clone(),
		Name:        name,
		IPVer:       ipVer,
		IP:          ip,
		Port:        port,
		ExitChan:    make(chan bool, 1),
		DataPack:    &msg.DefaultDataPack{},
		WorkerGroup: &impl.WorkerGroup{},
		ConnManager: &server2.ConnManager{
			ConnMap: make(map[uint32]iface.IConn),
		},
		HeaderOperator: &msg.DefaultHeaderOperator{},
	}
	dispatcher := &impl.DefaultDispatcher{
		ICatLog:       clog.AppLogger().Clone(),
		MsgHandlerMap: make(map[interface{}]func(ctx *iface.CatContext, request iface.IRequest, data []byte) error),
		Server:        ser,
	}
	ser.SetDispatcher(dispatcher)
	return ser
}
