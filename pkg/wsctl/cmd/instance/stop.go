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
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("accepts 1 arg(s), received %d", len(args))
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			if err := stop(cmd, client, args); err != nil {
				return errors.Wrap(err, fmt.Sprintf("problem stop instance %+v", args))
			}
			cmd.Println(cases.Title(language.Und).String(args[0]) + " instance stopped successfully ")
			return nil
		},
	}
}

func stop(cmd *cobra.Command, client client.Client, args []string) error {
	url := GetInstanceCmdUrl(client.Config().Endpoint, args[0], "STOP")
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return errors.Wrap(err, "failed to stop instance request")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Call(url, req)
	if err != nil {
		return errors.Wrap(err, "failed to stop instance")
	}
	return utils.PrintResponse(cmd, resp)
}
