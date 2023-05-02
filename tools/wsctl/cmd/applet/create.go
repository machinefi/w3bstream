package applet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/machinefi/w3bstream/tools/wsctl/client"
	"github.com/machinefi/w3bstream/tools/wsctl/config"
)

var (
	_appletCreateUse = map[config.Language]string{
		config.English: "create PROJECT_NAME FILE_PATH",
		config.Chinese: "create PROJECT_NAME FILE_PATH",
	}
	_appletCreateCmdShorts = map[config.Language]string{
		config.English: "Create a applet",
		config.Chinese: "创建 APPLET",
	}
)

// newAppletCreateCmd is a command to create applet
func newAppletCreateCmd(client client.Client) *cobra.Command {
	return &cobra.Command{
		Use:   client.SelectTranslation(_appletCreateUse),
		Short: client.SelectTranslation(_appletCreateCmdShorts),
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			if err := Create(client, args[0], args[1]); err != nil {
				return errors.Wrap(err, fmt.Sprintf("problem create applet %+v", args))
			}
			cmd.Println("applet created successfully")
			return nil
		},
	}
}

type info struct {
	AppletName string `json:"appletName,omitempty"`
	WasmName   string `json:"wasmName,omitempty"`
	WasmMd5    string `json:"wasmMd5,omitempty"`
}

func createURL(endpoint, name string) string {
	return fmt.Sprintf("%s/srv-applet-mgr/v0/project/x/%s", endpoint, name)
}

func Create(client client.Client, projectName, path string) error {
	body, err := prepareRequest(path)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", createURL(client.Config().Endpoint, projectName), body)
	if err != nil {
		return errors.Wrap(err, "failed to create applet request")
	}
	req.Header.Set("Content-Type", "application/json")

	_, err = client.Call(req)
	if err != nil {
		return errors.Wrap(err, "failed to create applet")
	}
	return nil
}

func prepareRequest(filePath string) (io.Reader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	filename := filepath.Base(file.Name())
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(part, file); err != nil {
		return nil, err
	}

	inf, err := json.Marshal(info{
		AppletName: "test",
		WasmName:   "test",
	})
	if err != nil {
		return nil, err
	}

	if err := writer.WriteField("info", string(inf)); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	return body, nil
}
