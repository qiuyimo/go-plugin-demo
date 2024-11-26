package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"demo/shared"
)

var (
	plugins map[string]*plugin.Client
	logger  hclog.Logger
)

const (
	ShadowPluginName = "shadow"
)

func initPlugin() {
	plugins[ShadowPluginName] = plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  shared.Handshake,
		Plugins:          shared.PluginMap,
		Cmd:              exec.Command("./plugin.exe"),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		Logger:           logger,
	})
}

func main() {
	logger = hclog.New(&hclog.LoggerOptions{
		Name:   "main",
		Output: os.Stdout,
		Level:  hclog.Info,
	})

	plugins = make(map[string]*plugin.Client)

	r := gin.Default()

	r.GET("/begin", func(c *gin.Context) {
		initPlugin()
		addr, err := plugins[ShadowPluginName].Start()
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("err: %v\n", err))
			return
		}
		c.String(http.StatusOK, "addr: %s", addr.String())
	})

	r.GET("/stop", func(c *gin.Context) {
		client, ok := plugins[ShadowPluginName]
		if !ok {
			c.String(http.StatusBadRequest, "shadow plugin not found")
			return
		}

		client.Kill()
		c.String(http.StatusOK, "stopped")
	})

	r.GET("/download", func(c *gin.Context) {
		client, ok := plugins[ShadowPluginName]
		if !ok {
			c.String(http.StatusBadRequest, "shadow plugin not found")
			return
		}

		rpcClient, err := client.Client()
		if err != nil {
			c.String(http.StatusBadRequest, "err: %v", err)
			return
		}

		raw, err := rpcClient.Dispense("shadow")
		if err != nil {
			c.String(http.StatusBadRequest, "err: %v", err)
			return
		}

		shadowCli := raw.(shared.ShadowInterface)
		res, err := shadowCli.Download("id:1", "name:a", "version:v1.1.1", "bucket:shadow")
		if err != nil {
			c.String(http.StatusBadRequest, "err: %v", err)
			return
		}

		c.String(http.StatusOK, "resp: %s", string(res))
	})

	r.Run(":8080")
}
