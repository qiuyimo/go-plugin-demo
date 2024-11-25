// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package shared

import (
	"context"

	"demo/proto"
)

// GRPCClient 主进程使用，业务接口KV的实现，通过 gRPC 客户端转发请求给插件进程
type GRPCClient struct{ client proto.ShadowClient }

func (m *GRPCClient) Download(key string, value []byte) error {
	// 将请求转发给插件进程
	_, err := m.client.Download(context.Background(), &proto.PutRequest{
		Key:   key,
		Value: value,
	})
	return err
}

// GRPCServer 插件进程使用
type GRPCServer struct {
	Impl ShadowInterface
	proto.UnimplementedShadowServer
}

func (m *GRPCServer) Download(ctx context.Context, req *proto.PutRequest) (*proto.Empty, error) {
	return &proto.Empty{}, m.Impl.Download(req.Key, req.Value)
}
