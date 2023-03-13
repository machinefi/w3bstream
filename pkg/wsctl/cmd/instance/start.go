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
	_instanceStartUse = map[config.Language]string{
		config.English: "start INSTANCE_ID",
		config.Chinese: "start INSTANCE_ID",
	}
	_instanceStartCmdShorts = map[config.Language]string{
		config.English: "Start a instance",
		config.Chinese: "通过 INSTANCE_ID 启动 INSTANCE",
	}
)

// newInstanceStartCmd is a command to start instance
func newInstanceStartCmd(client client.Client) *cobra.Command {
	return &cobra.Command{
		Use:   client.SelectTranslation(_instanceStartUse),
		Short: client.SelectTranslation(_instanceStartCmdShorts),
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("accepts 1 arg(s), received %d", len(args))
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			if err := start(cmd, client, args); err != nil {
				return errors.Wrap(err, fmt.Sprintf("problem start instance %+v", args))
			}
			cmd.Println(cases.Title(language.Und).String(args[0]) + " instance started successfully ")
			return nil
		},
	}
}

func start(cmd *cobra.Command, client client.Client, args []string) error {
	url := GetInstanceCmdUrl(client.Config().Endpoint, args[0], "START")
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return errors.Wrap(err, "failed to start instance request")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Call(url, req)
	if err != nil {
		return errors.Wrap(err, "failed to start instance")
	}
	return utils.PrintResponse(cmd, resp)
}
