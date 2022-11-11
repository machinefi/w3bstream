package main

import (
	"os"

	"github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/wsctl/client"
	"github.com/machinefi/w3bstream/pkg/wsctl/cmd"
	"github.com/machinefi/w3bstream/pkg/wsctl/cmd/config"
)

func main() {
	logger := log.Std()

	readConfig, defaultConfigFile, err := config.InitConfig()
	if err != nil {
		logger.Panic(err)
	}
	client := client.NewClient(readConfig, defaultConfigFile, logger)
	if err := cmd.NewWsctl(client).Execute(); err != nil {
		os.Exit(1)
	}
}
