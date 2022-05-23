package impl

import (
	"bytes"
	"encoding/binary"
)

// DefaultHeader 实现基本信息包头
type DefaultHeader struct {
	PackType uint32
	Data     []byte
}

func (bm *DefaultHeader) Pack() ([]byte, error) {
	buff := bytes.NewBuffer([]byte{})

	if err := binary.Write(buff, binary.LittleEndian, bm.PackType); err != nil {
		return nil, err
	}

	if err := binary.Write(buff, binary.LittleEndian, bm.Data); err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

func (bm *DefaultHeader) Unpack(data []byte) error {
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

func (bm *DefaultHeader) GetPackType() interface{} {
	return bm.PackType
}

func (bm *DefaultHeader) GetData() []byte {
	return bm.Data
}
