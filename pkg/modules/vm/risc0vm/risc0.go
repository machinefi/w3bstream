package risc0vm

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/gorilla/websocket"

	"github.com/machinefi/w3bstream/pkg/depends/kit/logr"
)

func CreateProof(ctx context.Context, req *CreateProofReq, host string, path string) (*CreateProofRsp, error) {
	params := make(map[string]interface{})

	ctx, l := logr.Start(ctx, "modules.xvm.CreateProof")
	defer l.End()

	params["params"] = req.Params
	params["image_id"] = req.ImageID
	jsonParams, err := json.Marshal(params)
	if err != nil {
		l.Error(err)
		return nil, err
	}

	url := url.URL{Scheme: "ws", Host: host, Path: path}
	client, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	defer client.Close()

	if err != nil {
		l.Error(err)
		return nil, err
	}
	err = client.WriteMessage(websocket.TextMessage, jsonParams)
	if err != nil {
		l.Error(err)
		return nil, err
	}

	_, message, err := client.ReadMessage()
	if err != nil {
		l.Error(err)
		return nil, err
	}
	println(string(message))

	return &CreateProofRsp{
		Receipt: string(message),
	}, nil
}
