package strategy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"

	"github.com/machinefi/w3bstream/tools/wsctl/client"
	"github.com/machinefi/w3bstream/tools/wsctl/config"
)

var (
	_strategyCreateUse = map[config.Language]string{
		config.English: "create PROJECT_NAME APPLET_ID EVENT_TYPE HANDLER",
		config.Chinese: "create PROJECT_NAME APPLET_ID EVENT_TYPE HANDLER",
	}
	_strategyCreateCmdShorts = map[config.Language]string{
		config.English: "Create a strategy",
		config.Chinese: "通过 PROJECT_ID, EVENT_TYPE HANDLER 创建 STRATEGY",
	}
)

// newStrategyCreateCmd is a command to create strategy
func newStrategyCreateCmd(client client.Client) *cobra.Command {
	return &cobra.Command{
		Use:   client.SelectTranslation(_strategyCreateUse),
		Short: client.SelectTranslation(_strategyCreateCmdShorts),
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			pubID, err := Create(client, args[0], args[1], args[2], args[3])
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("problem create strategy %+v", args))
			}
			cmd.Printf("strategy %s created successfully\n", pubID)
			return nil
		},
	}
}

type strategy struct {
	AppletID  string `json:"appletID"`
	EventType string `json:"eventType"`
	Handler   string `json:"handler"`
}

func createURL(endpoint, name string) string {
	return fmt.Sprintf("%s/srv-applet-mgr/v0/strategy/x/%s", endpoint, name)
}

func Create(client client.Client, projectname, appID, eventType, handler string) (string, error) {
	bodyBytes, err := json.Marshal(strategy{appID, eventType, handler})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", createURL(client.Config().Endpoint, projectname), bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", errors.Wrap(err, "failed to create strategy request")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Call(req)
	if err != nil {
		return "", errors.Wrap(err, "failed to create strategy")
	}
	if !gjson.ValidBytes(resp) {
		return "", errors.New("invalid response")
	}
	pubID := gjson.ParseBytes(resp).Get("strategyID")
	if !pubID.Exists() {
		return "", errors.New("invalid response")
	}
	return pubID.String(), nil
}
