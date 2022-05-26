package iface

import (
	"ccat/clog"
	"context"
)

type CatContext struct {
	context.Context
	clog.ICatLog
	IConn
}

func NewCatContext(ctx context.Context, conn IConn) *CatContext {
	return &CatContext{Context: ctx, ICatLog: clog.AppLogger().Clone(), IConn: conn}
}
