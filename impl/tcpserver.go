package impl

import (
	"ccat/iface"
	"fmt"
	"net"
)

// TcpService tcp 服务
type TcpService struct {
	Name       string
	IPVer      string
	IP         string
	Port       uint32
	DataPack   iface.IDataPack
	Dispatcher iface.IDispatcher
}

// Start 创建tcp监听
func (t *TcpService) Start() {
	lsn, err := net.Listen(t.IPVer, fmt.Sprintf("%s:%d", t.IP, t.Port))
	if err != nil {
		fmt.Println("[Tcp Service] listen err", err)
		return
	}

	fmt.Println("[Tcp Service] Start listen...")

	for {
		conn, err := lsn.Accept()
		if err != nil {
			fmt.Println("[Tcp Service] accept err", err)
			break
		}
		// 处理连接读
		catConn := &Conn{
			C:        conn,
			IsValid:  true,
			ExitChan: make(chan bool, 1),
			Server:   t,
		}
		go catConn.Start()
	}
}

// Stop 停止服务，关闭所有连接
func (t *TcpService) Stop() {

}

func (t *TcpService) Run() {
	t.Start()
}

// GetDataPack 获取数据包处理对象
func (t *TcpService) GetDataPack() iface.IDataPack {
	return t.DataPack
}

// SetDataPack 设置数据包处理对象
func (t *TcpService) SetDataPack(packDeal iface.IDataPack) {
	t.DataPack = packDeal
}

// SetDispatcher 设置消息分发器
func (t *TcpService) SetDispatcher(dispatcher iface.IDispatcher) {
	t.Dispatcher = dispatcher
}

// GetDispatcher 获取消息分发器
func (t *TcpService) GetDispatcher() iface.IDispatcher {
	return t.Dispatcher
}

// GetAddr 获取服务监听地址
func (t *TcpService) GetAddr() string {
	return fmt.Sprintf("%s:%d", t.IP, t.Port)
}

// GetName 获取服务名
func (t *TcpService) GetName() string {
	return t.Name
}
