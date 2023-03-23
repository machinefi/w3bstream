package publisher

import (
	"bytes"
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
	_publisherCreateUse = map[config.Language]string{
		config.English: "create PROJECT_ID PUB_NAME PUB_KEY",
		config.Chinese: "create PROJECT_ID PUB_NAME PUB_KEY",
	}
	_publisherCreateCmdShorts = map[config.Language]string{
		config.English: "Create a publisher",
		config.Chinese: "通过 PROJECT_ID, PUB_NAME, PUB_KEY 创建 PUBLISHER",
	}
)

// newPublisherCreateCmd is a command to create publisher
func newPublisherCreateCmd(client client.Client) *cobra.Command {
	return &cobra.Command{
		Use:   client.SelectTranslation(_publisherCreateUse),
		Short: client.SelectTranslation(_publisherCreateCmdShorts),
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 3 {
				return fmt.Errorf("accepts 3 arg(s), received %d", len(args))
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			if err := create(cmd, client, args); err != nil {
				return errors.Wrap(err, fmt.Sprintf("problem create publisher %+v", args))
			}
			cmd.Println(cases.Title(language.Und).String(args[0]) + " publisher created successfully ")
			return nil
		},
	}
}

func create(cmd *cobra.Command, client client.Client, args []string) error {
	body := fmt.Sprintf(`{"name":"%s", "key":"%s"}`, args[1], args[2])
	createURL := GetPublisherCmdUrl(client.Config().Endpoint, args[0])
	req, err := http.NewRequest("POST", createURL, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return errors.Wrap(err, "failed to create publisher request")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Call(createURL, req)
	if err != nil {
		return errors.Wrap(err, "failed to create publisher")
	}
	return utils.PrintResponse(cmd, resp)
}
