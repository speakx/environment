package srvinstance

import (
	"context"
	"errors"

	"google.golang.org/grpc"
)

// GrpcClient 封装
type GrpcClient struct {
	addr   string
	conn   *grpc.ClientConn
	ctx    context.Context
	cancel context.CancelFunc
}

// Connect connect
func (g *GrpcClient) Connect(addr string) error {
	if nil != g.conn {
		return errors.New("it's already connected")
	}
	if len(addr) == 0 {
		return errors.New("addr is nil")
	}

	// 开始正常连接操作
	g.addr = addr
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return err
	}

	g.conn = conn
	g.ctx, g.cancel = context.WithCancel(context.Background())
	return nil
}

// Disconnect disconnect
func (g *GrpcClient) Disconnect() error {
	if nil != g.conn {
		g.conn.Close()
		g.conn = nil
	}
	if nil != g.cancel {
		g.cancel()
		g.cancel = nil
	}
	return nil
}

// Reconnect disconnect + connect
func (g *GrpcClient) Reconnect() error {
	err := g.Disconnect()
	if nil != err {
		return err
	}
	return g.Connect(g.addr)
}

// GetConn conn
func (g *GrpcClient) GetConn() *grpc.ClientConn {
	return g.conn
}

// GetCtx ctx
func (g *GrpcClient) GetCtx() context.Context {
	return g.ctx
}
