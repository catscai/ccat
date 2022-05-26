package msg

import (
	"encoding/binary"
	"io"
	"net"
)

type DefaultDataPack struct {
}

var defaultHeaderLen = 4 // 默认头长度4个字节，表示包的长度

// ParseData 从连接中解析出包长度和数据
func (pack *DefaultDataPack) ParseData(conn net.Conn) ([]byte, error) {
	packLenBytes := make([]byte, defaultHeaderLen)
	if _, err := io.ReadFull(conn, packLenBytes); err != nil {
		return nil, err
	}
	// 解析出包长度
	packLen := binary.LittleEndian.Uint32(packLenBytes)

	var data []byte
	if packLen > 0 {
		data = make([]byte, packLen)
		if _, err := io.ReadFull(conn, data); err != nil {
			return nil, err
		}
	}

	return data, nil
}

// ReorganizeData 将消息数据重新组织为可发送的data
func (pack *DefaultDataPack) ReorganizeData(data []byte) ([]byte, error) {
	packLen := len(data)
	packLenBytes := make([]byte, defaultHeaderLen)
	binary.LittleEndian.PutUint32(packLenBytes, uint32(packLen))

	sendData := append(packLenBytes, data...)

	return sendData, nil
}
