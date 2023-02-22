package instance

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/machinefi/w3bstream/pkg/wsctl/client"
	"github.com/machinefi/w3bstream/pkg/wsctl/cmd/utils"
	"github.com/machinefi/w3bstream/pkg/wsctl/config"
)

var (
	_instanceDeleteUse = map[config.Language]string{
		config.English: "delete INSTANCE_ID",
		config.Chinese: "delete INSTANCE_ID",
	}
	_instanceDeleteCmdShorts = map[config.Language]string{
		config.English: "Delete a instance",
		config.Chinese: "通过 INSTANCE_ID 删除 INSTANCE",
	}
)

// newInstanceDeleteCmd is a command to delete instance
func newInstanceDeleteCmd(client client.Client) *cobra.Command {
	return &cobra.Command{
		Use:   client.SelectTranslation(_instanceDeleteUse),
		Short: client.SelectTranslation(_instanceDeleteCmdShorts),
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("accepts 1 arg(s), received %d", len(args))
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			if err := delete(cmd, client, args); err != nil {
				return errors.Wrap(err, fmt.Sprintf("problem delete instance %+v", args))
			}
			cmd.Println(cases.Title(language.Und).String(args[0]) + " instance deleted successfully ")
			return nil
		},
	}
}

func delete(cmd *cobra.Command, client client.Client, args []string) error {
	url := GetInstanceCmdUrl(client.Config().Endpoint, args[0], "REMOVE")
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return errors.Wrap(err, "failed to delete instance request")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Call(url, req)
	if err != nil {
		return errors.Wrap(err, "failed to delete instance")
	}
	return utils.PrintResponse(cmd, resp)
}
