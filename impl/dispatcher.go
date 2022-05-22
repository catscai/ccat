package impl

import (
	"ccat/iface"
	"fmt"
	"reflect"
)

// BaseDispatcher 保存包与回调业务映射关系，业务分发
type BaseDispatcher struct {
	MsgHandlerMap map[interface{}]func(conn iface.IConn, data []byte) error
}

// Dispatch 将消息分发给处理函数
func (bd *BaseDispatcher) Dispatch(conn iface.IConn, data []byte) {
	fmt.Println("Start Dispatch message")
	msg := DefaultMessage{}
	msg.Unpack(data)
	fmt.Printf("Dispatch Unpack msg:%+v\n", msg)
	if f, ok := bd.MsgHandlerMap[msg.GetPackType()]; ok {
		f(conn, msg.GetData())
	} else {
		fmt.Println("Not found message handler, packType", msg.GetPackType())
	}
}

// RegisterHandler 注册消息回调
func (bd *BaseDispatcher) RegisterHandler(packType interface{}, message iface.IMessage, deal iface.MsgHandlerFunc) {
	fmt.Println("RegisterHandler packType", packType)
	msgType := reflect.TypeOf(message).Elem()
	msgTypeName := msgType.String()
	fmt.Println("msgTypeName", msgTypeName)
	handler := func(conn iface.IConn, data []byte) error {
		// 利用反射创建新对象
		req := reflect.New(msgType).Elem().Addr().Interface().(iface.IMessage)
		fmt.Println("handler req", req)
		if err := req.Unpack(data); err != nil {
			fmt.Println("req Message Unpack err", err, "packName:", msgTypeName)
			return err
		}

		// todo 加一个recover panic 捕捉业务处理时(deal执行时)的异常情况

		// 调用业务回调
		if err := deal(conn, req); err != nil {
			fmt.Println("Business Deal err", err)
			return err
		}
		return nil
	}
	bd.MsgHandlerMap[packType] = handler
}

// Remove 删除回调映射关系
func (bd *BaseDispatcher) Remove(packType interface{}) {
	delete(bd.MsgHandlerMap, packType)
}
