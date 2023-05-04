package project

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/machinefi/w3bstream/tools/wsctl/client"
	"github.com/machinefi/w3bstream/tools/wsctl/config"
)

var (
	_projectDeleteUse = map[config.Language]string{
		config.English: "delete PROJECT_NAME",
		config.Chinese: "delete project名称",
	}
	_projectDeleteCmdShorts = map[config.Language]string{
		config.English: "Delete a project",
		config.Chinese: "删除一个project",
	}
)

// newProjectDeleteCmd is a command to delete project
func newProjectDeleteCmd(client client.Client) *cobra.Command {
	return &cobra.Command{
		Use:   client.SelectTranslation(_projectDeleteUse),
		Short: client.SelectTranslation(_projectDeleteCmdShorts),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			projectName := args[0]
			if err := Delete(client, projectName); err != nil {
				return errors.Wrap(err, fmt.Sprintf("problem delete project %+v", args))
			}
			cmd.Printf("project %s deleted successfully\n", projectName)
			return nil
		},
	}
}

func deleteURL(endpoint, name string) string {
	return fmt.Sprintf("%s/srv-applet-mgr/v0/project/x/%s", endpoint, name)
}

func Delete(client client.Client, name string) error {
	req, err := http.NewRequest("DELETE", deleteURL(client.Config().Endpoint, name), nil)
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}
	req.Header.Set("Content-Type", "application/json")
	if _, err = client.Call(req); err != nil {
		return err
	}
	return nil
}
