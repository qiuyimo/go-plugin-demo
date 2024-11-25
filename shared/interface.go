package shared

import (
	"context"

	"demo/proto"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

var PluginMap = map[string]plugin.Plugin{
	"shadow": &ShadowGRPCPlugin{},
}

// ShadowInterface 业务接口：这个是与 proto/kv.proto 保持一致的业务接口，注意要返回 error
type ShadowInterface interface {
	Download(id, name, version, bucket string) error
}

type ShadowGRPCPlugin struct {
	plugin.Plugin
	Impl ShadowInterface
}

func (p *ShadowGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	// 向grpc.ServerGRPC 类型参数s，注册服务的实现
	proto.RegisterShadowServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *ShadowGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	// 创建gRPC客户端的方法是自动生成的
	return &GRPCClient{client: proto.NewShadowClient(c)}, nil
}
