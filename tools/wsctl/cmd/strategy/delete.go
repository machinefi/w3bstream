package strategy

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/machinefi/w3bstream/tools/wsctl/client"
	"github.com/machinefi/w3bstream/tools/wsctl/config"
)

var (
	_strategyDeleteUse = map[config.Language]string{
		config.English: "delete STRATEGY_ID",
		config.Chinese: "delete STRATEGY_ID",
	}
	_strategyDeleteCmdShorts = map[config.Language]string{
		config.English: "Delete a strategy",
		config.Chinese: "删除 STRATEGY",
	}
)

// newStrategyDeleteCmd is a command to delete strategy
func newStrategyDeleteCmd(client client.Client) *cobra.Command {
	return &cobra.Command{
		Use:   client.SelectTranslation(_strategyDeleteUse),
		Short: client.SelectTranslation(_strategyDeleteCmdShorts),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			if err := Delete(client, args[0]); err != nil {
				return errors.Wrap(err, fmt.Sprintf("problem delete strategy %+v", args))
			}
			cmd.Printf("strategy %s deleted successfully\n", args[0])
			return nil
		},
	}
}

func deleteURL(endpoint, stgID string) string {
	return fmt.Sprintf("%s/srv-applet-mgr/v0/strategy/data/%s", endpoint, stgID)
}

func Delete(client client.Client, stgID string) error {
	req, err := http.NewRequest("DELETE", deleteURL(client.Config().Endpoint, stgID), nil)
	if err != nil {
		return errors.Wrap(err, "failed to delete strategy request")
	}
	req.Header.Set("Content-Type", "application/json")
	_, err = client.Call(req)
	if err != nil {
		return errors.Wrap(err, "failed to delete strategy")
	}
	return nil
}
