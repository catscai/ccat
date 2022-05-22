package impl

import (
	"bytes"
	"encoding/binary"
)

// DefaultMessage 实现基本的消息包
type DefaultMessage struct {
	PackType uint32
	Data     []byte
}

func (bm *DefaultMessage) Pack() ([]byte, error) {
	buff := bytes.NewBuffer([]byte{})

	if err := binary.Write(buff, binary.LittleEndian, bm.PackType); err != nil {
		return nil, err
	}

	if err := binary.Write(buff, binary.LittleEndian, bm.Data); err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

func (bm *DefaultMessage) Unpack(data []byte) error {
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

func (bm *DefaultMessage) GetPackType() interface{} {
	return bm.PackType
}

func (bm *DefaultMessage) GetData() []byte {
	return bm.Data
}
