package publisher

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
	_publisherDeleteUse = map[config.Language]string{
		config.English: "delete PROJECT_NAME",
		config.Chinese: "delete PROJECT_NAME",
	}
	_publisherDeleteCmdShorts = map[config.Language]string{
		config.English: "Delete a publisher",
		config.Chinese: "通过 PROJECT_NAME 删除 PUBLISHER",
	}
)

// newPublisherDeleteCmd is a command to delete publisher
func newPublisherDeleteCmd(client client.Client) *cobra.Command {
	return &cobra.Command{
		Use:   client.SelectTranslation(_publisherDeleteUse),
		Short: client.SelectTranslation(_publisherDeleteCmdShorts),
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("accepts 1 arg(s), received %d", len(args))
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			if err := delete(cmd, client, args); err != nil {
				return errors.Wrap(err, fmt.Sprintf("problem delete publisher %+v", args))
			}
			cmd.Println(cases.Title(language.Und).String(args[0]) + " publisher deleted successfully ")
			return nil
		},
	}
}

func delete(cmd *cobra.Command, client client.Client, args []string) error {
	url := GetPublisherCmdUrl(client.Config().Endpoint, args[0])
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return errors.Wrap(err, "failed to delete publisher request")
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Call(url, req)
	if err != nil {
		return errors.Wrap(err, "failed to delete publisher")
	}
	return utils.PrintResponse(cmd, resp)
}
