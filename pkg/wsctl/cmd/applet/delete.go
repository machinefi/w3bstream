package applet

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
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("accepts 1 arg(s), received %d", len(args))
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			if err := delete(cmd, client, args); err != nil {
				return errors.Wrap(err, fmt.Sprintf("problem delete applet %+v", args))
			}
			cmd.Println(cases.Title(language.Und).String(args[0]) + " applet deleted successfully ")
			return nil
		},
	}
}

func delete(cmd *cobra.Command, client client.Client, args []string) error {
	url := GetAppletCmdUrl(client.Config().Endpoint, args[0])
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return errors.Wrap(err, "failed to delete applet request")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Call(url, req)
	if err != nil {
		return errors.Wrap(err, "failed to delete applet")
	}
	return utils.PrintResponse(cmd, resp)
}
