package applet

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/machinefi/w3bstream/pkg/wsctl/client"
	"github.com/machinefi/w3bstream/pkg/wsctl/config"
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

func GetAppletCmdUrl(endpoint, cmd string) string {
	return fmt.Sprintf("%s/srv-applet-mgr/v0/applet/%s", endpoint, cmd)
}
