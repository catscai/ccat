package ccat

import (
	"ccat/config"
	"ccat/iface"
	"ccat/impl"
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
			AddServer(NewTcpService(name, "tcp4", info.IP, info.Port))
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
	return &impl.TcpService{
		Name:     name,
		IPVer:    ipVer,
		IP:       ip,
		Port:     port,
		ExitChan: make(chan bool, 1),
		DataPack: &impl.DefaultDataPack{},
		Dispatcher: &impl.BaseDispatcher{
			MsgHandlerMap: make(map[interface{}]func(conn iface.IConn, data []byte) error),
		},
		HeaderParser: &impl.DefaultHeaderParser{},
		WorkerGroup:  &impl.WorkerGroup{},
	}
}
