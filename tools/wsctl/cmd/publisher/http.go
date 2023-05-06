package publisher

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/machinefi/w3bstream/tools/wsctl/client"
	"github.com/machinefi/w3bstream/tools/wsctl/config"
)

var (
	_httpEventUse = map[config.Language]string{
		config.English: "http PROJECT_NAME TOKEN EVENT_TYPE PAYLOAD",
		config.Chinese: "http PROJECT_NAME TOKEN EVENT_TYPE PAYLOAD",
	}
	_httpEventCmdShorts = map[config.Language]string{
		config.English: "send an event via http protocol",
		config.Chinese: "send an event via http protocol",
	}
)

// newHttpEventCmd is a command to send event via http protocol
func newHttpEventCmd(client client.Client) *cobra.Command {
	return &cobra.Command{
		Use:   client.SelectTranslation(_httpEventUse),
		Short: client.SelectTranslation(_httpEventCmdShorts),
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			err := HTTPEvent(client, args[0], args[1], args[2], args[3])
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("problem sent event %+v", args))
			}
			cmd.Printf("message %s sent successfully\n", args[3])
			return nil
		},
	}
}

func httpURL(endpoint, name, eventType string) string {
	return fmt.Sprintf("http://localhost:8889/srv-applet-mgr/v0/event/%s?eventType=%s", name, eventType)
}

func HTTPEvent(client client.Client, projectname, token, eventType, payload string) error {
	req, err := http.NewRequest(
		"POST",
		httpURL(
			client.Config().Endpoint,
			fmt.Sprintf("eth_%s_%s", client.Address(), projectname),
			eventType,
		),
		bytes.NewBuffer([]byte(payload)),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create strategy request")
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	_, err = client.CallWithToken(req, token)
	if err != nil {
		return errors.Wrap(err, "failed to send event")
	}
	return nil
}
