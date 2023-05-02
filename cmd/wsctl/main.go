package main

import (
	"os"

	"github.com/machinefi/w3bstream/tools/wsctl/client"
	"github.com/machinefi/w3bstream/tools/wsctl/cmd"
	"github.com/machinefi/w3bstream/tools/wsctl/cmd/config"
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		panic(err)
	}
	client := client.NewClient(cfg)
	if err := cmd.NewWsctl(client).Execute(); err != nil {
		os.Exit(1)
	}
}
