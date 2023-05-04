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
	"github.com/tidwall/gjson"

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
			appID, _, err := Create(client, args[0], args[1])
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("problem create applet %+v", args))
			}
			cmd.Printf("applet %s created successfully\n", appID)
			return nil
		},
	}
}

type strategy struct {
	EventType string `json:"eventType"`
	Handler   string `json:"handler"`
}

type info struct {
	AppletName string `json:"appletName,omitempty"`
	Deploy     bool   `json:"start"`
	WasmName   string `json:"wasmName,omitempty"`
	// WasmMd5    string `json:"wasmMd5,omitempty"`
	// Strategies []strategy `json:"strategies,omitempty"`
}

func createURL(endpoint, name string) string {
	return fmt.Sprintf("%s/srv-applet-mgr/v0/applet/x/%s", endpoint, name)
}

func Create(client client.Client, projectName, path string) (string, string, error) {
	req, err := prepareRequest(path, createURL(client.Config().Endpoint, projectName))
	if err != nil {
		return "", "", err
	}
	resp, err := client.Call(req)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to create applet")
	}

	if !gjson.ValidBytes(resp) {
		return "", "", errors.New("invalid response")
	}
	appID := gjson.ParseBytes(resp).Get("instance.appletID")
	insID := gjson.ParseBytes(resp).Get("instance.instanceID")
	if !appID.Exists() || !insID.Exists() {
		return "", "", errors.New("invalid response")
	}
	return appID.String(), insID.String(), nil
}

func prepareRequest(filePath, url string) (*http.Request, error) {
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
		AppletName: filename,
		Deploy:     true,
		WasmName:   filename,
		// Strategies: []strategy{{
		// 	"DEFAULT", "start",
		// }},
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

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, nil
}
