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

func (s *ShadowImplement) Download(id, name, version, bucket string) ([]byte, error) {
	s.logger.Info("##############", "id", id, "name", name, "version", version, "bucket", bucket)
	return []byte("aaa"), nil
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
			"shadow_plugin": &shared.ShadowGRPCPlugin{Impl: shadowImpl},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
