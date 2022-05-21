package impl

import (
	"fmt"
	"net"
)

// TcpService tcp 服务
type TcpService struct {
	IPVer string
	IP    string
	Port  uint32
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
