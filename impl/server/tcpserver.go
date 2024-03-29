package server

import (
	"fmt"
	"github.com/catscai/ccat/clog"
	"github.com/catscai/ccat/config"
	"github.com/catscai/ccat/iface"
	"github.com/catscai/ccat/iface/imsg"
	"go.uber.org/zap"
	"net"
)

// TcpService tcp 服务
type TcpService struct {
	clog.ICatLog
	Name             string
	IPVer            string
	IP               string
	Port             uint32
	DataPack         imsg.IDataPack          // 处理tcp粘包
	Dispatcher       iface.IDispatcher       // 消息分发器
	WorkerGroup      iface.IWorkerGroup      // 工作者组
	ExitChan         chan bool               // 退出管道
	ConnManager      iface.IConnManager      // 连接管理器
	ConnStartHandler iface.ConnStatusHandler // 连接建立成功后的回调
	ConnEndHandler   iface.ConnStatusHandler // 连接关闭后的回调
	HeaderOperator   imsg.IHeaderOperator    // 包头操作者
}

// Start 创建tcp监听
func (t *TcpService) Start() {
	lsn, err := net.Listen(t.IPVer, fmt.Sprintf("%s:%d", t.IP, t.Port))
	if err != nil {
		t.Error("[Tcp Service] listen failed", zap.Any("err", err))
		return
	}

	t.Info("[Tcp Service] Start listen...")
	// 退出关闭 释放资源
	defer t.Stop()
	maxConnLimit := uint32(0)
	cfg := config.GetTcpServiceCfg(t.GetName())
	if cfg != nil {
		maxConnLimit = cfg.MaxConn
	}
	var connID uint32
	for {
		select {
		case <-t.ExitChan:
			return
		default:
			conn, err := lsn.Accept()
			if err != nil {
				t.Error("[Tcp Service] accept failed", zap.Any("err", err))
				continue
			}

			// 超出最大连接后，拒绝连接
			if maxConnLimit > 0 && connID >= maxConnLimit {
				t.Warn("[Tcp Service] the number of connection over limit",
					zap.Any("curNum", connID+1), zap.Any("maxConnLimit", maxConnLimit))
				conn.Close()
				continue
			}
			// 处理连接读
			catConn := NewConnection(conn, connID, t)

			// 将连接加入管理者
			t.ConnManager.Add(catConn)
			connID++
			go catConn.Start()
		}
	}
}

// Stop 停止服务，关闭所有连接
func (t *TcpService) Stop() {
	close(t.ExitChan)
	t.WorkerGroup.Stop()
	t.ConnManager.Clear()
}

func (t *TcpService) Run() {
	cfg := config.GetTcpServiceCfg(t.GetName())
	t.WorkerGroup.Init(t.ICatLog, t, cfg.WorkerGroup.Size, cfg.WorkerGroup.QueueLength)
	t.WorkerGroup.Start()

	go t.Start()

	t.Info("[Tcp Service] Running...")
}

// GetDataPack 获取数据包处理对象
func (t *TcpService) GetDataPack() imsg.IDataPack {
	return t.DataPack
}

// SetDataPack 设置数据包处理对象
func (t *TcpService) SetDataPack(packDeal imsg.IDataPack) {
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

// GetWorkerGroup 获取工作者组
func (t *TcpService) GetWorkerGroup() iface.IWorkerGroup {
	return t.WorkerGroup
}

// GetConnManager 获取连接管理器
func (t *TcpService) GetConnManager() iface.IConnManager {
	return t.ConnManager
}

// SetConnectStartHandler 设置连接建立成功后的回调
func (t *TcpService) SetConnectStartHandler(handler iface.ConnStatusHandler) {
	t.ConnStartHandler = handler
}

// SetConnectEndHandler 设置连接关闭时的回调
func (t *TcpService) SetConnectEndHandler(handler iface.ConnStatusHandler) {
	t.ConnEndHandler = handler
}

// CallConnectStart 调用连接开始回调
func (t *TcpService) CallConnectStart(conn iface.IConn) {
	if t.ConnStartHandler != nil {
		t.ConnStartHandler(t, conn)
	}
}

// CallConnectEnd 调用连接退出回调
func (t *TcpService) CallConnectEnd(conn iface.IConn) {
	if t.ConnEndHandler != nil {
		t.ConnEndHandler(t, conn)
	}
}

// SetHeaderOperator 设置包头操作对象
func (t *TcpService) SetHeaderOperator(operator imsg.IHeaderOperator) {
	t.HeaderOperator = operator
}

// GetHeaderOperator 获取包头操作对象
func (t *TcpService) GetHeaderOperator() imsg.IHeaderOperator {
	return t.HeaderOperator
}

// GetAddr 获取服务监听地址
func (t *TcpService) GetAddr() string {
	return fmt.Sprintf("%s:%d", t.IP, t.Port)
}

// GetName 获取服务名
func (t *TcpService) GetName() string {
	return t.Name
}
