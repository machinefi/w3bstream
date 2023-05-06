package applet

import (
	"github.com/spf13/cobra"

	"github.com/machinefi/w3bstream/tools/wsctl/client"
	"github.com/machinefi/w3bstream/tools/wsctl/config"
)

// Multi-language support
var (
	_appletCmdShorts = map[config.Language]string{
		config.English: "Manage applets of W3bstream",
		config.Chinese: "管理 W3bstream 系统里的 applets",
	}
)

// NewAppletCmd represents the new applet command.
func NewAppletCmd(client client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "applet",
		Short: client.SelectTranslation(_appletCmdShorts),
	}
	cmd.AddCommand(newAppletDeleteCmd(client))
	cmd.AddCommand(newAppletCreateCmd(client))
	return cmd
}
