package client

import (
	"ccat/clog"
	"ccat/iface/imsg"
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net"
	"sync"
	"time"
)

type Client struct {
	clog.ICatLog
	Conn           net.Conn
	DataPack       imsg.IDataPack
	HeaderOperator imsg.IHeaderOperator
	process        func(conn net.Conn, header imsg.IHeaderPack) error
	isValid        bool
	sendQueue      chan imsg.IHeaderPack
	sessionChanMap map[interface{}]chan []byte
	mutex          sync.RWMutex
	timeOut        time.Duration
	ctx            context.Context
	cancel         context.CancelFunc
}

func NewClient(logger clog.ICatLog, dataPack imsg.IDataPack, headerOperator imsg.IHeaderOperator, sendChanLen uint32, sendTimeOut time.Duration) *Client {
	client := &Client{
		ICatLog:        logger,
		DataPack:       dataPack,
		HeaderOperator: headerOperator,
		sendQueue:      make(chan imsg.IHeaderPack, sendChanLen),
		sessionChanMap: make(map[interface{}]chan []byte),
		isValid:        false,
		timeOut:        sendTimeOut,
		process:        nil,
	}
	return client
}

// Connection 连接服务器
func (client *Client) Connection(ipVer, address string, timeout time.Duration) error {
	conn, err := net.DialTimeout(ipVer, address, timeout)
	if err != nil {
		client.Error(fmt.Sprintf("[Client] Connection failed, ipVer:%s,address:%s",
			ipVer, address), zap.Any("error", err))
		return err
	}
	client.Conn = conn
	client.isValid = true

	client.ctx, client.cancel = context.WithCancel(context.Background())
	// 连接成功,创建读写协程
	go client.beginRead()
	go client.beginWrite()
	go client.release()
	return nil
}

// SetProcess 设置消息回调
func (client *Client) SetProcess(process func(conn net.Conn, header imsg.IHeaderPack) error) {
	client.process = process
}

// SetDataPack 设置处理粘包，分包
func (client *Client) SetDataPack(pack imsg.IDataPack) {
	client.DataPack = pack
}

// SetHeaderOperator 设置包头解析
func (client *Client) SetHeaderOperator(operator imsg.IHeaderOperator) {
	client.HeaderOperator = operator
}

// Send 同步发送，等待请求回复;如果设置了process,那么收到消息后既会调用process回调，也会在Send接口返回
func (client *Client) Send(req, rsp imsg.IHeaderPack) error {
	if !client.Valid() {
		return fmt.Errorf("invalid client")
	}
	resChan, err := client.setChan(req.GetSessionID())
	if err != nil {
		return err
	}
	defer client.delChan(req.GetSessionID())

	client.sendQueue <- req
	t := time.NewTimer(client.timeOut)
	select {
	case <-t.C:
		client.Error("[Client] Send timeout...", zap.Any("req", req))
		return errors.New("client send timeout")
	case data := <-resChan:
		if err = rsp.Unpack(data); err != nil {
			client.Error("[Client] rsp.Unpack failed", zap.Any("error", err))
			return errors.New("response msg Unpack err")
		}
		client.Info("[Client] Success Recv data", zap.Any("session", rsp.GetSessionID()))
	}
	return nil
}

// SendASync 异步发送
func (client *Client) SendASync(req imsg.IHeaderPack) error {
	if client.Valid() {
		client.sendQueue <- req
		return nil
	}

	return errors.New("client is invalid")
}

// Valid 连接是否有效
func (client *Client) Valid() bool {
	return client.isValid
}

// beginRead 连接是否有效
func (client *Client) beginRead() {
	client.Info("[Client] beginRead start...")
	defer client.Close()
	for {
		select {
		case <-client.ctx.Done():
			return
		default:
			data, err := client.DataPack.ParseData(client.Conn)
			if err != nil {
				client.Error("[Client] DataPack.ParseData failed", zap.Any("error", err))
				return
			}
			header := client.HeaderOperator.Get()
			if err = header.Unpack(data); err != nil {
				client.Error("[Client] HeaderParser.HeaderUnpack failed", zap.Any("error", err))
				return
			}

			resChan := client.getChan(header.GetSessionID())
			if resChan != nil {
				resChan <- data
			} else {
				client.Warn("client.getChan resChan is nil", zap.Any("session", header.GetSessionID()))
			}
			if client.process != nil {
				client.process(client.Conn, header)
			}
		}
	}
}

func (client *Client) beginWrite() {
	client.Info("[Client] beginWrite start...")
	defer client.Close()
	for {
		select {
		case <-client.ctx.Done():
			return
		case header := <-client.sendQueue:
			// 发送队列已经是用户封装好的header了，所以不需要再次封装包头
			// client.HeaderParser.HeaderPack(packType,header.GetData())
			data, err := header.Pack()
			if err != nil {
				client.Error("[Client] beginWrite header.Pack failed",
					zap.Any("error", err), zap.Any("session", header.GetSessionID()))
				continue
			}
			sendData, err := client.DataPack.ReorganizeData(data)
			if err != nil {
				client.Error("[Client] beginWrite DataPack.ReorganizeData failed", zap.Any("error", err))
				continue
			}
			// 阻塞写则一定是全部发送的，不需要关心发送了多少
			if _, err = client.Conn.Write(sendData); err != nil {
				client.Error("[Client] beginWrite Conn.Write failed", zap.Any("error", err))
				return
			}
		}
	}
}

// Close 关闭连接
func (client *Client) Close() {
	client.cancel() // 发送取消信号
}

func (client *Client) setChan(sessionID interface{}) (chan []byte, error) {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	if _, ok := client.sessionChanMap[sessionID]; ok {
		return nil, errors.New("chan already exist")
	}
	client.sessionChanMap[sessionID] = make(chan []byte, 1)

	return client.sessionChanMap[sessionID], nil
}

func (client *Client) getChan(sessionID interface{}) chan []byte {
	client.mutex.RLock()
	defer client.mutex.RUnlock()
	if c, ok := client.sessionChanMap[sessionID]; ok {
		return c
	}
	return nil
}

func (client *Client) delChan(sessionID interface{}) {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	if _, ok := client.sessionChanMap[sessionID]; ok {
		close(client.sessionChanMap[sessionID])
		delete(client.sessionChanMap, sessionID)
	}
}

func (client *Client) release() {
	select {
	case <-client.ctx.Done():
		client.Info("[Client] release connection close", zap.Any("remoteAddr", client.Conn.RemoteAddr()))
		client.Conn.Close()
		client.isValid = false
		close(client.sendQueue)
		return
	}
}
