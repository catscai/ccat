package impl

import (
	"ccat/iface"
	"ccat/iface/imsg"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

type Client struct {
	Conn           net.Conn
	DataPack       iface.IDataPack
	HeaderParser   iface.IHeaderPackParser
	process        func(conn net.Conn, header imsg.IHeaderPack) error
	isValid        bool
	exitChan       chan bool
	sendQueue      chan imsg.IHeaderPack
	sessionChanMap map[interface{}]chan []byte
	mutex          sync.RWMutex
	timeOut        uint32
}

// Connection 连接服务器
func (client *Client) Connection(ipVer, address string, chanLen, timeout uint32) error {
	fmt.Println("[Client] Connection start...", "ipVer", ipVer, "address", address)
	conn, err := net.Dial(ipVer, address)
	if err != nil {
		fmt.Println("[Client] Connection err", err)
		return err
	}
	client.Conn = conn
	client.isValid = true
	client.sendQueue = make(chan imsg.IHeaderPack, chanLen)
	client.timeOut = timeout

	// 连接成功,创建读写协程
	go client.beginRead()
	go client.beginWrite()
	return nil
}

// SetProcess 设置消息回调
func (client *Client) SetProcess(process func(conn net.Conn, header imsg.IHeaderPack) error) {
	client.process = process
}

// SetDataPack 设置处理粘包，分包
func (client *Client) SetDataPack(pack iface.IDataPack) {
	client.DataPack = pack
}

// SetHeaderParser 设置包头解析
func (client *Client) SetHeaderParser(parser iface.IHeaderPackParser) {
	client.HeaderParser = parser
}

// Send 同步发送，等待请求回复;如果设置了process,那么收到消息后既会调用process回调，也会在Send接口返回
func (client *Client) Send(req, rsp imsg.IHeaderPack) error {
	resChan, err := client.setChan(req.GetSessionID())
	if err != nil {
		return err
	}
	defer client.delChan(req.GetSessionID())

	client.sendQueue <- req
	t := time.NewTimer(time.Millisecond * time.Duration(client.timeOut))
	select {
	case <-t.C:
		fmt.Println("[Client] Send timeout...")
		return errors.New("client send timeout")
	case data := <-resChan:
		if err = rsp.Unpack(data); err != nil {
			fmt.Println("[Client] rsp.Unpack err", err)
			return errors.New("response msg Unpack err")
		}
		fmt.Println("[Client] Success Recv data")
	}
	return nil
}

// SendASync 异步发送
func (client *Client) SendASync(req imsg.IHeaderPack) error {
	client.sendQueue <- req
	return nil
}

// Valid 连接是否有效
func (client *Client) Valid() bool {
	return client.isValid
}

// beginRead 连接是否有效
func (client *Client) beginRead() {
	fmt.Println("[Client] beginRead start...")
	defer client.Close()
	for {
		select {
		case <-client.exitChan:
			return
		default:
			data, err := client.DataPack.ParseData(client.Conn)
			if err != nil {
				fmt.Println("[Client] DataPack.ParseData err", err)
				return
			}
			header, err := client.HeaderParser.HeaderUnpack(data)
			if err != nil {
				fmt.Println("[Client] HeaderParser.HeaderUnpack err", err)
				return
			}

			fmt.Println("[Client] Recv data", header)
			if client.process != nil {
				client.process(client.Conn, header)
			}
			resChan := client.getChan(header.GetSessionID())
			if resChan != nil {
				resChan <- data
			}
		}
	}
}

func (client *Client) beginWrite() {
	fmt.Println("[Client] beginWrite start...")
	defer client.Close()

	for {
		select {
		case <-client.exitChan:
			return
		case header := <-client.sendQueue:
			// 发送队列已经是用户封装好的header了，所以不需要再次封装包头
			// client.HeaderParser.HeaderPack(req.GetPackType(),req.GetData())
			data, err := header.Pack()
			if err != nil {
				fmt.Println("[Client] beginWrite  header.Pack err", err, "header", header)
				continue
			}
			sendData, err := client.DataPack.ReorganizeData(data)
			if err != nil {
				fmt.Println("[Client] beginWrite DataPack.ReorganizeData err", err)
				continue
			}
			// todo 确认知识点 golang默认阻塞写？ 若是阻塞写则一定是全部发送的，不需要关心发送了多小
			if _, err = client.Conn.Write(sendData); err != nil {
				fmt.Println("[Client] beginWrite Conn.Write err", err)
				return
			}
		}
	}
}

// Close 关闭连接
func (client *Client) Close() {
	client.Conn.Close()
	client.isValid = false
	close(client.exitChan)
	close(client.sendQueue)
	fmt.Println("[Client] Close..")
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
