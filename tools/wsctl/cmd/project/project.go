package project

import (
	"github.com/spf13/cobra"

	"github.com/machinefi/w3bstream/tools/wsctl/client"
	"github.com/machinefi/w3bstream/tools/wsctl/config"
)

// Multi-language support
var (
	_projectCmdShorts = map[config.Language]string{
		config.English: "Manage projects of W3bstream",
		config.Chinese: "管理 W3bstream 系统里的 projects",
	}
)

// NewProjectCmd represents the new project command.
func NewProjectCmd(client client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: client.SelectTranslation(_projectCmdShorts),
	}
	cmd.AddCommand(newProjectCreateCmd(client))
	cmd.AddCommand(newProjectDeleteCmd(client))
	return cmd
}
