package publisher

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/machinefi/w3bstream/tools/wsctl/client"
	"github.com/machinefi/w3bstream/tools/wsctl/config"
)

var (
	_publisherDeleteUse = map[config.Language]string{
		config.English: "delete PUBLISHER_ID",
		config.Chinese: "delete PUBLISHER_ID",
	}
	_publisherDeleteCmdShorts = map[config.Language]string{
		config.English: "Delete a publisher",
		config.Chinese: "删除 PUBLISHER",
	}
)

// newPublisherDeleteCmd is a command to delete publisher
func newPublisherDeleteCmd(client client.Client) *cobra.Command {
	return &cobra.Command{
		Use:   client.SelectTranslation(_publisherDeleteUse),
		Short: client.SelectTranslation(_publisherDeleteCmdShorts),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			if err := Delete(client, args[0]); err != nil {
				return errors.Wrap(err, fmt.Sprintf("problem delete publisher %+v", args))
			}
			cmd.Printf("publisher %s deleted successfully\n", args[0])
			return nil
		},
	}
}

func deleteURL(endpoint, pubID string) string {
	return fmt.Sprintf("%s/srv-applet-mgr/v0/publisher/data/%s", endpoint, pubID)
}

func Delete(client client.Client, pubID string) error {
	req, err := http.NewRequest("DELETE", deleteURL(client.Config().Endpoint, pubID), nil)
	if err != nil {
		return errors.Wrap(err, "failed to delete publisher request")
	}
	req.Header.Set("Content-Type", "application/json")
	_, err = client.Call(req)
	if err != nil {
		return errors.Wrap(err, "failed to delete publisher")
	}
	return nil
}
