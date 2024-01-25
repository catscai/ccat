package iface

import (
	"github.com/catscai/ccat/iface/imsg"
	"github.com/golang/protobuf/proto"
)

// IDispatcher 分发消息
type IDispatcher interface {
	Dispatch(request IRequest)                                                        // 将消息分发给处理函数
	RegisterHandler(packType interface{}, message imsg.IMessage, deal MsgHandlerFunc) // 注册消息回调
	// RegisterHandlerSimple 简单回调
	RegisterHandlerSimple(reqType, rspType interface{}, reqMsg, rspMsg imsg.IMessage, deal MsgHandlerSimpleFunc)
	RegisterHandlerPB(reqType interface{}, message proto.Message, deal MsgHandlerFuncPB)
	RegisterHandlerSimplePB(reqType, rspType interface{}, reqMsg, rspMsg proto.Message, deal MsgHandlerSimpleFuncPB)

	RegisterHandlerData(reqType interface{}, message imsg.IHeaderPack, deal MsgHandlerFuncData)
	Remove(reqType interface{}) // 删除回调映射关系
}

type MsgHandlerFunc func(ctx *CatContext, request IRequest, iMessage imsg.IMessage) error
type MsgHandlerSimpleFunc func(ctx *CatContext, reqMsg, rspMsg imsg.IMessage) error

type MsgHandlerFuncPB func(ctx *CatContext, request IRequest, iMessage proto.Message) error
type MsgHandlerSimpleFuncPB func(ctx *CatContext, reqMsg, rspMsg proto.Message) error

type MsgHandlerFuncData func(ctx *CatContext, request IRequest, iMessage imsg.IHeaderPack) error
