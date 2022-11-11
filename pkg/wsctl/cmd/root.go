package cmd

import (
	"github.com/spf13/cobra"

	"github.com/machinefi/w3bstream/pkg/wsctl/client"
	cfgcmd "github.com/machinefi/w3bstream/pkg/wsctl/cmd/config"
	"github.com/machinefi/w3bstream/pkg/wsctl/config"
)

// Multi-language support
var (
	_wsctlRootCmdShorts = map[config.Language]string{
		config.English: "Command-line interface for Machinefi W3bstream",
		config.Chinese: "Machinefi W3bstream 命令行工具",
	}
	_wsctlRootCmdLongs = map[config.Language]string{
		config.English: `wsctl is a command-line interface for interacting with Machinefi W3bstream`,
		config.Chinese: `wsctl 是用于与 Machinefi W3bstream 进行交互的命令行工具`,
	}
)

// NewWsctl returns wsctl root cmd
func NewWsctl(client client.Client) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "wsctl",
		Short: client.SelectTranslation(_wsctlRootCmdShorts),
		Long:  client.SelectTranslation(_wsctlRootCmdLongs),
	}

	rootCmd.AddCommand(cfgcmd.NewConfigCmd(client))

	return rootCmd
}
