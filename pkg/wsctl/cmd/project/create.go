package project

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/machinefi/w3bstream/pkg/wsctl/client"
	"github.com/machinefi/w3bstream/pkg/wsctl/config"
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
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("accepts 1 arg(s), received %d", len(args))
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			result, err := create(client, args)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("problem create project %+v", args))
			}
			cmd.Println(result)
			return nil
		},
	}
}

type projectResp struct {
	Code uint64 `json:"code"`
	Desc string `json:"desc"`
}

func create(client client.Client, args []string) (string, error) {
	body := fmt.Sprintf(`{"name":"%s"}`, args[0])
	url := fmt.Sprintf("%s/srv-applet-mgr/v0/project", client.Config().Endpoint)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return "", errors.Wrap(err, "failed to create project request")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Call(url, req)
	if err != nil {
		return "", errors.Wrap(err, "failed to create project")
	}
	defer resp.Body.Close()

	cr := projectResp{}
	if err := json.NewDecoder(resp.Body).Decode(&cr); err != nil {
		errors.Wrap(err, "failed to decode project responce")
	}
	if cr.Code != 0 {
		return "", fmt.Errorf("failed to create project, error code: %d, error message: %s", cr.Code, cr.Desc)
	}

	return cases.Title(language.Und).String(args[0]) + " project created successfully ", nil
}
