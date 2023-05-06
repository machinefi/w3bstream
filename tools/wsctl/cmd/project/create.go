package project

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/machinefi/w3bstream/tools/wsctl/client"
	"github.com/machinefi/w3bstream/tools/wsctl/config"
)

var (
	_projectCreateUse = map[config.Language]string{
		config.English: "create PROJECT_NAME",
		config.Chinese: "create project名称",
	}
	_projectCreateCmdShorts = map[config.Language]string{
		config.English: "Create a new project",
		config.Chinese: "创建一个新的project",
	}
)

// newProjectCreateCmd is a command to create project
func newProjectCreateCmd(client client.Client) *cobra.Command {
	return &cobra.Command{
		Use:   client.SelectTranslation(_projectCreateUse),
		Short: client.SelectTranslation(_projectCreateCmdShorts),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			projectName := args[0]
			if _, err := Create(client, projectName); err != nil {
				return errors.Wrap(err, fmt.Sprintf("failed to create project %s", projectName))
			}
			cmd.Printf("project %s created successfully\n", projectName)
			return nil
		},
	}
}

type projectCreate struct {
	Name string `json:"name"`
}

func createURL(endpoint string) string {
	return fmt.Sprintf("%s/srv-applet-mgr/v0/project", endpoint)
}

func Create(client client.Client, name string) ([]byte, error) {
	bodyBytes, err := json.Marshal(projectCreate{name})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(
		"POST",
		createURL(client.Config().Endpoint),
		bytes.NewBuffer(bodyBytes),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}
	req.Header.Set("Content-Type", "application/json")
	return client.Call(req)
}
