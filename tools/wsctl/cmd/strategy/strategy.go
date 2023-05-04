package strategy

import (
	"github.com/spf13/cobra"

	"github.com/machinefi/w3bstream/tools/wsctl/client"
	"github.com/machinefi/w3bstream/tools/wsctl/config"
)

// Multi-language support
var (
	_strategyCmdShorts = map[config.Language]string{
		config.English: "Manage strategies of W3bstream",
		config.Chinese: "管理 W3bstream 系统里的 strategies",
	}
)

// NewStrategyCmd represents the new strategy command.
func NewStrategyCmd(client client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "strategy",
		Short: client.SelectTranslation(_strategyCmdShorts),
	}
	cmd.AddCommand(newStrategyDeleteCmd(client))
	cmd.AddCommand(newStrategyCreateCmd(client))
	return cmd
}
