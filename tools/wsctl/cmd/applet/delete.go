package applet

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/machinefi/w3bstream/tools/wsctl/client"
	"github.com/machinefi/w3bstream/tools/wsctl/config"
)

var (
	_appletDeleteUse = map[config.Language]string{
		config.English: "delete APPLET_ID",
		config.Chinese: "delete APPLET_ID",
	}
	_appletDeleteCmdShorts = map[config.Language]string{
		config.English: "Delete a applet",
		config.Chinese: "通过 APPLET_ID 删除 APPLET",
	}
)

// newAppletDeleteCmd is a command to delete applet
func newAppletDeleteCmd(client client.Client) *cobra.Command {
	return &cobra.Command{
		Use:   client.SelectTranslation(_appletDeleteUse),
		Short: client.SelectTranslation(_appletDeleteCmdShorts),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			if err := Delete(client, args[0]); err != nil {
				return errors.Wrap(err, fmt.Sprintf("problem delete applet %+v", args))
			}
			cmd.Printf("applet %s deleted successfully\n", args[0])
			return nil
		},
	}
}

func deleteURL(endpoint, appID string) string {
	return fmt.Sprintf("%s/srv-applet-mgr/v0/applet/data/%s", endpoint, appID)
}

func Delete(client client.Client, appID string) error {
	req, err := http.NewRequest("DELETE", deleteURL(client.Config().Endpoint, appID), nil)
	if err != nil {
		return errors.Wrap(err, "failed to delete applet request")
	}
	req.Header.Set("Content-Type", "application/json")

	_, err = client.Call(req)
	if err != nil {
		return errors.Wrap(err, "failed to delete applet")
	}
	return nil
}
