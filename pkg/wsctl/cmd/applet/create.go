package applet

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/machinefi/w3bstream/pkg/wsctl/client"
	"github.com/machinefi/w3bstream/pkg/wsctl/cmd/utils"
	"github.com/machinefi/w3bstream/pkg/wsctl/config"
)

var (
	_appletCreateUse = map[config.Language]string{
		config.English: "create PROJECT_ID FILE INFO",
		config.Chinese: "create PROJECT_ID FILE INFO",
	}
	_appletCreateCmdShorts = map[config.Language]string{
		config.English: "Create a applet",
		config.Chinese: "通过 PROJECT_ID, FILE, INFO 创建 APPLET",
	}
)

// newAppletCreateCmd is a command to create applet
func newAppletCreateCmd(client client.Client) *cobra.Command {
	return &cobra.Command{
		Use:   client.SelectTranslation(_appletCreateUse),
		Short: client.SelectTranslation(_appletCreateCmdShorts),
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 3 {
				return fmt.Errorf("accepts 3 arg(s), received %d", len(args))
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			if err := create(cmd, client, args); err != nil {
				return errors.Wrap(err, fmt.Sprintf("problem create applet %+v", args))
			}
			cmd.Println(cases.Title(language.Und).String(args[0]) + " applet created successfully ")
			return nil
		},
	}
}

func create(cmd *cobra.Command, client client.Client, args []string) error {
	createURL := GetAppletCmdUrl(client.Config().Endpoint, args[0])
	body, err := loadFile(args[1], args[2])
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", createURL, body)
	if err != nil {
		return errors.Wrap(err, "failed to create applet request")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Call(createURL, req)
	if err != nil {
		return errors.Wrap(err, "failed to create applet")
	}
	return utils.PrintResponse(cmd, resp)
}

func loadFile(filePath string, info string) (io.Reader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(part, file); err != nil {
		return nil, err
	}

	if err := writer.WriteField("info", info); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	return body, nil
}
