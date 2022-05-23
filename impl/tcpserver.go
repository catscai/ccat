package impl

import (
	"ccat/config"
	"ccat/iface"
	"fmt"
	"net"
)

// TcpService tcp 服务
type TcpService struct {
	Name         string
	IPVer        string
	IP           string
	Port         uint32
	DataPack     iface.IDataPack         // 处理tcp粘包
	Dispatcher   iface.IDispatcher       // 消息分发器
	HeaderParser iface.IHeaderPackParser // 包头解析器
	WorkerGroup  iface.IWorkerGroup      // 工作者组
	ExitChan     chan bool               // 退出管道
}

// Start 创建tcp监听
func (t *TcpService) Start() {
	lsn, err := net.Listen(t.IPVer, fmt.Sprintf("%s:%d", t.IP, t.Port))
	if err != nil {
		fmt.Println("[Tcp Service] listen err", err)
		return
	}

	fmt.Println("[Tcp Service] Start listen...")
	// 退出关闭 释放资源
	defer t.Stop()

	for {
		select {
		case <-t.ExitChan:
			return
		default:
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
}

// Stop 停止服务，关闭所有连接
func (t *TcpService) Stop() {
	close(t.ExitChan)
	t.WorkerGroup.Stop()
}

func (t *TcpService) Run() {
	t.DataPack.Init(t)
	cfg := config.GetTcpServiceCfg(t.GetName())
	t.WorkerGroup.Init(t, cfg.WorkerGroup.Size, cfg.WorkerGroup.QueueLength)
	t.WorkerGroup.Start()

	go t.Start()

	fmt.Println("[Tcp Service] Running...")
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

// SetHeaderParser 设置包头解析器
func (t *TcpService) SetHeaderParser(parser iface.IHeaderPackParser) {
	t.HeaderParser = parser
}

// GetHeaderParser 获取包头解析器
func (t *TcpService) GetHeaderParser() iface.IHeaderPackParser {
	return t.HeaderParser
}

// GetWorkerGroup 获取工作者组
func (t *TcpService) GetWorkerGroup() iface.IWorkerGroup {
	return t.WorkerGroup
}

// GetAddr 获取服务监听地址
func (t *TcpService) GetAddr() string {
	return fmt.Sprintf("%s:%d", t.IP, t.Port)
}

// GetName 获取服务名
func (t *TcpService) GetName() string {
	return t.Name
}
