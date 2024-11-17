package main

import (
    "fmt"
    "os"
    
    "github.com/hashicorp/go-plugin"
    
    "demo/shared"
)

// KVImplement 业务接口的实现
type KVImplement struct{}

func (KVImplement) Put(key string, value []byte) error {
    value = []byte(fmt.Sprintf("%s\n\nWritten from plugin-go-grpc", string(value)))
    return os.WriteFile("kv_"+key, value, 0644)
}

func (KVImplement) Get(key string) ([]byte, error) {
    return os.ReadFile("kv_" + key)
}

func main() {
    plugin.Serve(&plugin.ServeConfig{
        HandshakeConfig: shared.Handshake,
        Plugins: map[string]plugin.Plugin{
            "kv": &shared.KVGRPCPlugin{Impl: &KVImplement{}},
        },
        
        // A non-nil value here enables gRPC serving for this plugin...
        GRPCServer: plugin.DefaultGRPCServer,
    })
}
