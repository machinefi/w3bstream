package publisher

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/machinefi/w3bstream/pkg/wsctl/client"
	"github.com/machinefi/w3bstream/pkg/wsctl/config"
)

// Multi-language support
var (
	_publisherCmdShorts = map[config.Language]string{
		config.English: "Manage publishers of W3bstream",
		config.Chinese: "管理 W3bstream 系统里的 publishers",
	}
)

// NewPublisherCmd represents the new publisher command.
func NewPublisherCmd(client client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "publisher",
		Short: client.SelectTranslation(_publisherCmdShorts),
	}
	cmd.AddCommand(newPublisherDeleteCmd(client))
	cmd.AddCommand(newPublisherCreateCmd(client))
	return cmd
}

func GetPublisherCmdUrl(endpoint, cmd string) string {
	return fmt.Sprintf("%s/srv-publisher-mgr/v0/publisher/%s", endpoint, cmd)
}
