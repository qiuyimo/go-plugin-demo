package main

import (
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"demo/shared"
)

// ShadowImplement 业务接口的实现
type ShadowImplement struct {
	logger hclog.Logger
}

func (s *ShadowImplement) Download(key string, value []byte) error {
	s.logger.Debug("##############key: %v, value: %v", key, string(value))
	return nil
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	shadowImpl := &ShadowImplement{
		logger: logger,
	}

	logger.Debug("message from plugin", "foo", "bar")

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"kv": &shared.KVGRPCPlugin{Impl: shadowImpl},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
