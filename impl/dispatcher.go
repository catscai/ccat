package impl

import (
	"ccat/iface"
	"ccat/iface/imsg"
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"reflect"
)

// DefaultDispatcher 保存包与回调业务映射关系，业务分发
type DefaultDispatcher struct {
	MsgHandlerMap map[interface{}]func(ctx *iface.CatContext, request iface.IRequest, data []byte) error
	Server        iface.IServer
}

// Dispatch 将消息分发给处理函数
func (bd *DefaultDispatcher) Dispatch(request iface.IRequest) {
	if f, ok := bd.MsgHandlerMap[request.GetHeaderPack().GetPackType()]; ok {
		f(iface.NewCatContext(context.TODO(), request.GetConn()), request,
			request.GetHeaderPack().GetData())
	} else {
		fmt.Println("Not found message handler, packType", request.GetHeaderPack().GetPackType())
	}
}

// RegisterHandler 注册消息回调
func (bd *DefaultDispatcher) RegisterHandler(packType interface{}, message imsg.IMessage, deal iface.MsgHandlerFunc) {
	fmt.Println("RegisterHandler packType", packType)
	msgType := reflect.TypeOf(message).Elem()
	msgTypeName := msgType.String()
	fmt.Println("msgTypeName", msgTypeName)
	handler := func(ctx *iface.CatContext, request iface.IRequest, data []byte) error {
		// 利用反射创建新对象
		req := reflect.New(msgType).Elem().Addr().Interface().(imsg.IMessage)
		fmt.Println("handler req", req)
		if err := req.Unpack(data); err != nil {
			fmt.Println("req Message Unpack err", err, "packName:", msgTypeName)
			return err
		}

		// todo 加一个recover panic 捕捉业务处理时(deal执行时)的异常情况
		defer RecoverPanic()
		// 调用业务回调
		if err := deal(ctx, request, req); err != nil {
			fmt.Println("Business Deal err", err)
			return err
		}
		return nil
	}
	bd.MsgHandlerMap[packType] = handler
}

// RegisterHandlerSimple 一次交互,自动发送
func (bd *DefaultDispatcher) RegisterHandlerSimple(reqType, rspType interface{},
	reqMsg, rspMsg imsg.IMessage, deal iface.MsgHandlerSimpleFunc) {
	reqMsgType := reflect.TypeOf(reqMsg).Elem()
	rspMsgType := reflect.TypeOf(rspMsg).Elem()

	handler := func(ctx *iface.CatContext, request iface.IRequest, data []byte) error {
		req := reflect.New(reqMsgType).Elem().Addr().Interface().(imsg.IMessage)
		rsp := reflect.New(rspMsgType).Elem().Addr().Interface().(imsg.IMessage)
		if err := req.Unpack(data); err != nil {
			fmt.Println("req Message Unpack err", err, "packName:", reqMsgType.String())
			return err
		}
		defer RecoverPanic()

		defer func() {
			rspData, err := rsp.Pack()
			if err != nil {
				fmt.Println("[DefaultDispatcher] rsp.Pack err", err)
				return
			}
			pkg := bd.Server.GetHeaderOperator().Full(rspType, rspData, request.GetHeaderPack())
			if err := request.GetConn().SendMsg(pkg); err != nil {
				fmt.Println("SendMsg err", err)
			}
		}()
		// 调用业务回调
		if err := deal(ctx, req, rsp); err != nil {
			fmt.Println("Business Deal err", err)
		}

		return nil
	}
	bd.MsgHandlerMap[reqType] = handler
}

// RegisterHandlerPB 用户自己控制发送 pb
func (bd *DefaultDispatcher) RegisterHandlerPB(reqType interface{}, message proto.Message, deal iface.MsgHandlerFuncPB) {
	reqMsgType := reflect.TypeOf(message).Elem()
	handler := func(ctx *iface.CatContext, request iface.IRequest, data []byte) error {
		req := reflect.New(reqMsgType).Elem().Addr().Interface().(proto.Message)
		if err := proto.Unmarshal(data, req); err != nil {
			fmt.Println("req Message Unpack err", err, "packName:", reqMsgType.String())
			return err
		}
		defer RecoverPanic()

		// 调用业务回调
		if err := deal(ctx, request, req); err != nil {
			fmt.Println("Business Deal err", err)
		}
		return nil
	}
	bd.MsgHandlerMap[reqType] = handler
}

// RegisterHandlerSimplePB 一次交互,自动发送,pb
func (bd *DefaultDispatcher) RegisterHandlerSimplePB(reqType, rspType interface{}, reqMsg, rspMsg proto.Message, deal iface.MsgHandlerSimpleFuncPB) {
	reqMsgType := reflect.TypeOf(reqMsg).Elem()
	rspMsgType := reflect.TypeOf(rspMsg).Elem()
	handler := func(ctx *iface.CatContext, request iface.IRequest, data []byte) error {
		req := reflect.New(reqMsgType).Elem().Addr().Interface().(proto.Message)
		rsp := reflect.New(rspMsgType).Elem().Addr().Interface().(proto.Message)
		if err := proto.Unmarshal(data, req); err != nil {
			fmt.Println("req Message Unpack err", err, "packName:", reqMsgType.String())
			return err
		}
		defer RecoverPanic()
		defer func() {
			rspData, err := proto.Marshal(rsp)
			if err != nil {
				fmt.Println("[DefaultDispatcher] proto.Marshal err", err)
				return
			}
			pkg := bd.Server.GetHeaderOperator().Full(rspType, rspData, request.GetHeaderPack())
			if err = request.GetConn().SendMsg(pkg); err != nil {
				fmt.Println("[DefaultDispatcher] SendMsg err", err)
			}
		}()
		// 调用业务回调
		if err := deal(ctx, req, rsp); err != nil {
			fmt.Println("Business Deal err", err)
		}

		return nil
	}
	bd.MsgHandlerMap[reqType] = handler
}

// RegisterHandlerData 注册回调 返回原始数据包,交给用户自己解析
func (bd *DefaultDispatcher) RegisterHandlerData(reqType interface{}, message imsg.IHeaderPack, deal iface.MsgHandlerFuncData) {
	handler := func(ctx *iface.CatContext, request iface.IRequest, data []byte) error {
		if err := deal(ctx, request, request.GetHeaderPack()); err != nil {
			fmt.Println("[DefaultDispatcher] deal err", err)
		}

		return nil
	}
	bd.MsgHandlerMap[reqType] = handler
}

// Remove 删除回调映射关系
func (bd *DefaultDispatcher) Remove(packType interface{}) {
	delete(bd.MsgHandlerMap, packType)
}

func RecoverPanic() {
	if r := recover(); r != nil {
		fmt.Println("[Panic] deal req:", r)
	}
}
