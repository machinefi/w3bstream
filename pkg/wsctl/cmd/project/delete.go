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
	_projectDeleteUse = map[config.Language]string{
		config.English: "delete PROJECT_NAME",
		config.Chinese: "delete 项目名称",
	}
	_projectDeleteCmdShorts = map[config.Language]string{
		config.English: "Delete a project",
		config.Chinese: "删除一个项目",
	}
)

// newProjectDeleteCmd is a command to delete project
func newProjectDeleteCmd(client client.Client) *cobra.Command {
	return &cobra.Command{
		Use:   client.SelectTranslation(_projectDeleteUse),
		Short: client.SelectTranslation(_projectDeleteCmdShorts),
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("accepts 1 arg(s), received %d", len(args))
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			result, err := delete(client, args)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("problem when deleting project %+v", args))
			}
			cmd.Println(result)
			return nil
		},
	}
}

func delete(client client.Client, args []string) (string, error) {
	url := fmt.Sprintf("%s/srv-applet-mgr/v0/project/%s", client.Config().Endpoint, args[0])
	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer([]byte{}))
	if err != nil {
		return "", errors.Wrap(err, "failed while sending delete project request")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Call(url, req)
	if err != nil {
		return "", errors.Wrap(err, "failed to delete project")
	}
	defer resp.Body.Close()

	cr := projectResp{}
	if err := json.NewDecoder(resp.Body).Decode(&cr); err != nil {
		errors.Wrap(err, "failed to decode project responce")
	}
	if cr.Code != 0 {
		return "", fmt.Errorf("failed to delete project, error code: %d, error message: %s", cr.Code, cr.Desc)
	}

	return cases.Title(language.Und).String(args[0]) + " project deleted successfully ", nil
}
