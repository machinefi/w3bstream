package project

import (
	"github.com/machinefi/w3bstream/pkg/wsctl/client"
	"github.com/machinefi/w3bstream/pkg/wsctl/config"
	"github.com/spf13/cobra"
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

	return cmd
}
