package instance

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/machinefi/w3bstream/pkg/wsctl/client"
	"github.com/machinefi/w3bstream/pkg/wsctl/config"
)

// Multi-language support
var (
	_instanceCmdShorts = map[config.Language]string{
		config.English: "Manage instances of W3bstream",
		config.Chinese: "管理 W3bstream 系统里的 instances",
	}
)

// NewInstanceCmd represents the new instance command.
func NewInstanceCmd(client client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "instance",
		Short: client.SelectTranslation(_instanceCmdShorts),
	}
	cmd.AddCommand(newInstanceStartCmd(client))
	cmd.AddCommand(newInstanceStopCmd(client))
	cmd.AddCommand(newInstanceDeleteCmd(client))

	return cmd
}

func GetInstanceCmdUrl(endpoint, insId, cmd string) string {
	return fmt.Sprintf("%s/srv-applet-mgr/v0/deploy/%s/%s", endpoint, insId, cmd)
}
