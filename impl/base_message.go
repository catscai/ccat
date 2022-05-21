package impl

import (
	"bytes"
	"encoding/binary"
)

// BaseMessage 实现基本的消息包
type BaseMessage struct {
	PackType uint32
	Data     []byte
}

func (bm *BaseMessage) Pack() ([]byte, error) {
	buff := bytes.NewBuffer([]byte{})

	if err := binary.Write(buff, binary.LittleEndian, bm.PackType); err != nil {
		return nil, err
	}

	if err := binary.Write(buff, binary.LittleEndian, bm.Data); err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

func (bm *BaseMessage) Unpack(data []byte) error {
	reader := bytes.NewReader(data)

	if err := binary.Read(reader, binary.LittleEndian, &bm.PackType); err != nil {
		return err
	}

	//if err := binary.Read(reader, binary.LittleEndian, &bm.Data); err != nil {
	//	return err
	//}

	bm.Data = data[4:]
	return nil
}
