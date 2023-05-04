package e2e

import (
	"os"
	"testing"

	"github.com/machinefi/w3bstream/tools/wsctl/client"
	"github.com/machinefi/w3bstream/tools/wsctl/cmd/config"
)

var (
	_cli client.Client
)

func TestMain(m *testing.M) {
	cfg, err := config.InitConfig()
	if err != nil {
		panic(err)
	}
	_cli = client.NewClient(cfg)

	exitVal := m.Run()

	os.Exit(exitVal)
}
