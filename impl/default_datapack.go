package impl

import (
	"ccat/config"
	"ccat/iface"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type DefaultDataPack struct {
	Server iface.IServer
}

var defaultHeaderLen = 4 // 默认头长度4个字节，表示包的长度

// Init 初始化，设置所属server
func (pack *DefaultDataPack) Init(server iface.IServer) {
	pack.Server = server
}

// ParseData 从连接中解析出包长度和数据
func (pack *DefaultDataPack) ParseData(conn iface.IConn) ([]byte, error) {
	packLenBytes := make([]byte, defaultHeaderLen)
	if _, err := io.ReadFull(conn.GetConn(), packLenBytes); err != nil {
		fmt.Println("ParseData ReadFull err", err)
		return nil, err
	}
	// 解析出包长度
	packLen := binary.LittleEndian.Uint32(packLenBytes)
	cfg := config.GetBaseServiceCfg(pack.Server.GetName())
	if cfg.MaxPackLen > 0 && packLen > cfg.MaxPackLen {
		fmt.Println("ParseData Recv packlen over max pack length limit, packLen", packLen)
		return nil, errors.New("packet length over max limit")
	}
	var data []byte
	if packLen > 0 {
		data = make([]byte, packLen)
		if _, err := io.ReadFull(conn.GetConn(), data); err != nil {
			fmt.Println("ParseData ReadFull err", err)
			return nil, err
		}
	}

	return data, nil
}

// ReorganizeData 将消息数据重新组织为可发送的data
func (pack *DefaultDataPack) ReorganizeData(msg iface.IMessage) ([]byte, error) {
	packData, err := msg.Pack()
	if err != nil {
		fmt.Println("ReorganizeData msg Pack err", err)
		return nil, err
	}
	packLen := len(packData)
	packLenBytes := make([]byte, defaultHeaderLen)
	binary.LittleEndian.PutUint32(packLenBytes, uint32(packLen))

	data := append(packLenBytes, packData...)

	return data, nil
}
