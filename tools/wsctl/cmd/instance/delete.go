package instance

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/machinefi/w3bstream/tools/wsctl/client"
	"github.com/machinefi/w3bstream/tools/wsctl/config"
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
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			if err := Delete(client, args[0]); err != nil {
				return errors.Wrap(err, fmt.Sprintf("problem delete instance %+v", args))
			}
			cmd.Printf("instance %s deleted successfully\n", args[0])
			return nil
		},
	}
}

func Delete(client client.Client, insID string) error {
	url := getInstanceCmdUrl(client.Config().Endpoint, insID, "KILL")
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return errors.Wrap(err, "failed to delete instance request")
	}
	req.Header.Set("Content-Type", "application/json")

	_, err = client.Call(req)
	if err != nil {
		return errors.Wrap(err, "failed to delete instance")
	}
	return nil
}
