package project

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/machinefi/w3bstream/tools/wsctl/client"
	"github.com/machinefi/w3bstream/tools/wsctl/cmd/utils"
	"github.com/machinefi/w3bstream/tools/wsctl/config"
)

var (
	_projectDeleteUse = map[config.Language]string{
		config.English: "delete PROJECT_NAME",
		config.Chinese: "delete project名称",
	}
	_projectDeleteCmdShorts = map[config.Language]string{
		config.English: "Delete a new project",
		config.Chinese: "删除一个新的project",
	}
)

// newProjectDeleteCmd is a command to delete project
func newProjectDeleteCmd(client client.Client) *cobra.Command {
	return &cobra.Command{
		Use:   client.SelectTranslation(_projectDeleteUse),
		Short: client.SelectTranslation(_projectDeleteCmdShorts),
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("accepts 1 arg(s), received %d", len(args))
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			if err := delete(cmd, client, args); err != nil {
				return errors.Wrap(err, fmt.Sprintf("problem delete project %+v", args))
			}
			cmd.Println(cases.Title(language.Und).String(args[0]) + " project deleted successfully ")
			return nil
		},
	}
}

func delete(cmd *cobra.Command, client client.Client, args []string) error {
	url := fmt.Sprintf("%s/srv-applet-mgr/v0/project/%s", client.Config().Endpoint, args[0])
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return errors.Wrap(err, "failed to delete project request")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Call(url, req)
	if err != nil {
		return errors.Wrap(err, "failed to delete project")
	}
	return utils.PrintResponse(cmd, resp)
}
