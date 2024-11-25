// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package shared

import (
	"context"

	"demo/proto"
)

// GRPCClient 主进程使用，业务接口KV的实现，通过 gRPC 客户端转发请求给插件进程
type GRPCClient struct{ client proto.ShadowClient }

// Download todo(rain): 之前一直报错，就是这里的返回只返回了 error，需要实现业务接口
func (m *GRPCClient) Download(id, name, version, bucket string) ([]byte, error) {
	// 将请求转发给插件进程
	resp, err := m.client.Download(context.Background(), &proto.DownloadReq{
		Id:      id,
		Name:    name,
		Version: version,
		Bucket:  bucket,
	})
	if err != nil {
		return nil, err
	}

	return resp.Value, err
}

// GRPCServer 插件进程使用
type GRPCServer struct {
	Impl ShadowInterface
	proto.UnimplementedShadowServer
}

func (m *GRPCServer) Download(ctx context.Context, req *proto.DownloadReq) (*proto.GetResponse, error) {
	b, err := m.Impl.Download(req.Id, req.Name, req.Version, req.Bucket)
	return &proto.GetResponse{Value: b}, err
}
