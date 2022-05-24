package ccat

import (
	"ccat/config"
	"ccat/iface"
	"ccat/impl"
	"ccat/impl/msg"
	server2 "ccat/impl/server"
)

// 全局项目app
var gApp *impl.App

func init() {
	gApp = &impl.App{
		ServerMap: make(map[string]iface.IServer),
	}
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
	return &server2.TcpService{
		Name:     name,
		IPVer:    ipVer,
		IP:       ip,
		Port:     port,
		ExitChan: make(chan bool, 1),
		DataPack: &msg.DefaultDataPack{},
		Dispatcher: &impl.DefaultDispatcher{
			MsgHandlerMap: make(map[interface{}]func(request iface.IRequest, data []byte) error),
		},
		HeaderParser: &msg.DefaultHeaderParser{},
		WorkerGroup:  &impl.WorkerGroup{},
		ConnManager: &server2.ConnManager{
			ConnMap: make(map[uint32]iface.IConn),
		},
	}
}
