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
	_instanceStopUse = map[config.Language]string{
		config.English: "stop INSTANCE_ID",
		config.Chinese: "stop INSTANCE_ID",
	}
	_instanceStopCmdShorts = map[config.Language]string{
		config.English: "Stop a instance",
		config.Chinese: "通过 INSTANCE_ID 停止 INSTANCE",
	}
)

// newInstanceStopCmd is a command to stop instance
func newInstanceStopCmd(client client.Client) *cobra.Command {
	return &cobra.Command{
		Use:   client.SelectTranslation(_instanceStopUse),
		Short: client.SelectTranslation(_instanceStopCmdShorts),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			if err := Stop(client, args[0]); err != nil {
				return errors.Wrap(err, fmt.Sprintf("problem stop instance %+v", args))
			}
			cmd.Printf("instance %s stopped successfully\n", args[0])
			return nil
		},
	}
}

func Stop(client client.Client, insID string) error {
	url := getInstanceCmdUrl(client.Config().Endpoint, insID, "HUNGUP")
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return errors.Wrap(err, "failed to stop instance request")
	}
	req.Header.Set("Content-Type", "application/json")

	_, err = client.Call(req)
	if err != nil {
		return errors.Wrap(err, "failed to stop instance")
	}
	return nil
}
