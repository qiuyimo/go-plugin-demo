// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package shared

import (
    "context"
    
    "demo/proto"
)

// GRPCClient 主进程使用，业务接口KV的实现，通过 gRPC 客户端转发请求给插件进程
type GRPCClient struct{ client proto.KVClient }

func (m *GRPCClient) Put(key string, value []byte) error {
    // 将请求转发给插件进程
    _, err := m.client.Put(context.Background(), &proto.PutRequest{
        Key:   key,
        Value: value,
    })
    return err
}

func (m *GRPCClient) Get(key string) ([]byte, error) {
    // 将请求转发给插件进程
    resp, err := m.client.Get(context.Background(), &proto.GetRequest{
        Key: key,
    })
    if err != nil {
        return nil, err
    }
    
    return resp.Value, nil
}

// GRPCServer 插件进程使用
type GRPCServer struct {
    Impl KV
    proto.UnimplementedKVServer
}

func (m *GRPCServer) Put(ctx context.Context, req *proto.PutRequest) (*proto.Empty, error) {
    return &proto.Empty{}, m.Impl.Put(req.Key, req.Value)
}

func (m *GRPCServer) Get(ctx context.Context, req *proto.GetRequest) (*proto.GetResponse, error) {
    v, err := m.Impl.Get(req.Key)
    return &proto.GetResponse{Value: v}, err
}
