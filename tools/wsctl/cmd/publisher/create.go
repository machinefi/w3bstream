package publisher

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
	_publisherCreateUse = map[config.Language]string{
		config.English: "create PROJECT_NAME PUB_NAME PUB_KEY",
		config.Chinese: "create PROJECT_NAME PUB_NAME PUB_KEY",
	}
	_publisherCreateCmdShorts = map[config.Language]string{
		config.English: "Create a publisher",
		config.Chinese: "通过 PROJECT_NAME, PUB_NAME, PUB_KEY 创建 PUBLISHER",
	}
)

// newPublisherCreateCmd is a command to create publisher
func newPublisherCreateCmd(client client.Client) *cobra.Command {
	return &cobra.Command{
		Use:   client.SelectTranslation(_publisherCreateUse),
		Short: client.SelectTranslation(_publisherCreateCmdShorts),
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			pubID, token, err := Create(client, args[0], args[1], args[2])
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("problem create publisher %+v", args))
			}
			cmd.Printf("publisher %s created successfully, publisher's token %s\n ", pubID, token)
			return nil
		},
	}
}

type publisher struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

func createURL(endpoint, name string) string {
	return fmt.Sprintf("%s/srv-applet-mgr/v0/publisher/x/%s", endpoint, name)
}

func Create(client client.Client, projectname, name, key string) (string, string, error) {
	bodyBytes, err := json.Marshal(publisher{name, key})
	if err != nil {
		return "", "", err
	}
	req, err := http.NewRequest("POST", createURL(client.Config().Endpoint, projectname), bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", "", errors.Wrap(err, "failed to create publisher request")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Call(req)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to create publisher")
	}
	if !gjson.ValidBytes(resp) {
		return "", "", errors.New("invalid response")
	}
	pubID := gjson.ParseBytes(resp).Get("publisherID")
	pubToken := gjson.ParseBytes(resp).Get("token")
	if !pubID.Exists() || !pubToken.Exists() {
		return "", "", errors.New("invalid response")
	}
	return pubID.String(), pubToken.String(), nil
}
