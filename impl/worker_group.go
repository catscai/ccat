package impl

import (
	"ccat/clog"
	"ccat/iface"
	"fmt"
	"go.uber.org/zap"
	"time"
)

type WorkerGroup struct {
	WorkerQueues []chan iface.IRequest        // 工作任务队列
	ExitChan     chan bool                    // 退出消息管道
	Server       iface.IServer                // 所属server
	ShardFunc    iface.ShardWorkerHandlerFunc // 用户可自定义，根据请求选择工作者处理
	logger       clog.ICatLog                 // 日志对象
}

// Init 初始化工作者组参数
func (wg *WorkerGroup) Init(logger clog.ICatLog, server iface.IServer, size uint32, queueLen uint32) {
	wg.Server = server
	wg.ExitChan = make(chan bool, 1)
	wg.WorkerQueues = make([]chan iface.IRequest, size)
	wg.logger = logger
	for i := 0; i < int(size); i++ {
		wg.WorkerQueues[i] = make(chan iface.IRequest, queueLen)
	}
}

// Start 启动工作者组
func (wg *WorkerGroup) Start() {
	size := wg.GetWorkerSize()
	for i := 0; i < int(size); i++ {
		go wg.work(uint32(i), wg.WorkerQueues[i])
	}
}

// Stop 停止，释放资源
func (wg *WorkerGroup) Stop() {
	close(wg.ExitChan)
}

// SendTask 添加处理请求任务
func (wg *WorkerGroup) SendTask(request iface.IRequest) {
	// 根据请求获取对应的worker id
	var id uint32
	if wg.ShardFunc == nil {
		id = wg.DefaultShardFunc(request)
	} else {
		id = wg.ShardFunc(wg, request)
	}
	if id >= wg.GetWorkerSize() {
		wg.logger.Error("[WorkerGroup] SendTask select worker id over max length", zap.Uint32("id", id))
		return
	}
	wg.WorkerQueues[id] <- request
}

// GetWorkerSize 获取工作者数量，工作者数量可在配置文件中配置
func (wg *WorkerGroup) GetWorkerSize() uint32 {
	return uint32(len(wg.WorkerQueues))
}

func (wg *WorkerGroup) work(id uint32, q chan iface.IRequest) {
	defer close(q)
	defer fmt.Printf("[WorkerGroup] id = %d exit", id)
	for {
		select {
		case <-wg.ExitChan:
			return
		case req := <-q:
			// 分发处理消息
			wg.Server.GetDispatcher().Dispatch(req)
		}
	}

}

// SetShardWorkerFunc 设置工作者选择回调,用户提供回调计算,该请求应该被分给哪个工作者执行
func (wg *WorkerGroup) SetShardWorkerFunc(f iface.ShardWorkerHandlerFunc) {
	wg.ShardFunc = f
}

// DefaultShardFunc 默认的工作者选择
func (wg *WorkerGroup) DefaultShardFunc(request iface.IRequest) uint32 {
	workerSize := wg.GetWorkerSize()
	now := time.Now().UnixNano()

	return uint32(now % int64(workerSize))
}
