package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"demo/shared"
)

func run() error {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "main",
		Output: os.Stdout,
		Level:  hclog.Info,
	})

	// We're a host. Start by launching the plugin process.
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  shared.Handshake,
		Plugins:          shared.PluginMap,
		Cmd:              exec.Command("./plugin.exe"),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		Logger:           logger,
	})
	defer client.Kill()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		return err
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("shadow")
	if err != nil {
		return err
	}

	// We should have a ShadowInterface store now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	shadowCli := raw.(shared.ShadowInterface)
	m, err := shadowCli.Download("id:1", "name:a", "version:v1.1.1", "bucket:shadow")
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	logger.Info("main complete", "m", m)

	return nil
}

func main() {
	// We don't want to see the plugin logs.
	log.SetOutput(ioutil.Discard)

	if err := run(); err != nil {
		fmt.Printf("error: %+v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}
