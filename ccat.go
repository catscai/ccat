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
		Servers: make([]iface.IServer, 0),
	}
	if config.AppCfg.IsTcpService {
		tcpCfg := &config.AppCfg.TcpCfg
		AddServer(NewTcpService("tcp4", tcpCfg.IP, tcpCfg.Port))
	}
}

func AddServer(server iface.IServer) {
	gApp.AddServer(server)
}

func Run() {
	gApp.Run()
}

func NewTcpService(IPVer, IP string, port uint32) iface.IServer {
	return &impl.TcpService{
		IPVer: IPVer,
		IP:    IP,
		Port:  port,
	}
}
