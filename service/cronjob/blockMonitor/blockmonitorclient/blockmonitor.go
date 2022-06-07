// Code generated by goctl. DO NOT EDIT!
// Source: governanceMonitor.proto

package blockmonitorclient

import (
	"context"
	blockmonitor2 "github.com/zecrey-labs/zecrey-legend/service/cronjob/blockMonitor/blockMonitor"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	Request  = blockmonitor2.Request
	Response = blockmonitor2.Response

	BlockMonitor interface {
		Ping(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
	}

	defaultBlockMonitor struct {
		cli zrpc.Client
	}
)

func NewBlockMonitor(cli zrpc.Client) BlockMonitor {
	return &defaultBlockMonitor{
		cli: cli,
	}
}

func (m *defaultBlockMonitor) Ping(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	client := blockmonitor2.NewBlockMonitorClient(m.cli.Conn())
	return client.Ping(ctx, in, opts...)
}
