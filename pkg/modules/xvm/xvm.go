package xvm

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/kit/logr"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

func CreateRisc0Vm(ctx context.Context, req *CreateRisc0VmReq) (*CreateRisc0VmRsp, error) {
	var (
		b bytes.Buffer

		prj = types.MustProjectFromContext(ctx)
		w   = multipart.NewWriter(&b)
	)

	_ = w.WriteField("project_name", prj.Name)

	fields := map[string]string{
		"image_id": req.ImageId,
		"elf":      req.Elf,
	}
	for key, val := range fields {
		_ = w.WriteField(key, val)
	}

	fw, err := w.CreateFormFile("file", "risc0.rs")
	if err != nil {
		panic(err)
	}
	fd, err := req.File.Open()
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(fw, fd)
	if err != nil {
		panic(err)
	}
	err = w.Close()
	if err != nil {
		panic(err)
	}
	request, err := http.NewRequest("POST", "http://127.0.0.1:3000/api/create_instance", &b)
	if err != nil {
		panic(err)
	}
	request.Header.Set("Content-Type", w.FormDataContentType())
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	return nil, nil
}

func CreateProof(ctx context.Context, req *CreateProofReq) (*CreateProofRsp, error) {
	var (
		prj = types.MustProjectFromContext(ctx)
		idg = confid.MustSFIDGeneratorFromContext(ctx)
		d   = types.MustMgrDBExecutorFromContext(ctx)

		params map[string]interface{}
	)

	ctx, l := logr.Start(ctx, "modules.xvm.CreateProof")
	defer l.End()

	proof := &models.Proof{
		RelProject: models.RelProject{ProjectID: prj.ProjectID},
		RelProof:   models.RelProof{ProofID: idg.MustGenSFID()},
		ProofInfo: models.ProofInfo{
			Name:         req.Name,
			TemplateName: req.TemplateName,
			ImageID:      req.ImageId,
			InputData:    req.InputData,
			Status:       enums.PROOF_STATUS__GENERATING,
		},
	}

	if err := json.Unmarshal([]byte(proof.InputData), &params); err != nil {
		l.Error(err)
		return nil, err
	}
	params["image_id"] = proof.ImageID
	jsonParams, err := json.Marshal(params)
	if err != nil {
		l.Error(err)
		return nil, err
	}

	if err := proof.Create(d); err != nil {
		if sqlx.DBErr(err).IsConflict() {
			return nil, status.ProofConflict
		}
		return nil, status.DatabaseError.StatusErr().WithDesc(err.Error())
	}

	url := url.URL{Scheme: "ws", Host: "127.0.0.1:3000", Path: "/ws/api/prove_file"}
	client, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		l.Error(err)
		return nil, err
	}
	err = client.WriteMessage(websocket.TextMessage, jsonParams)
	if err != nil {
		l.Error(err)
		return nil, err
	}

	go func() {
		defer client.Close()
		_, message, err := client.ReadMessage()
		if err != nil {
			l.Error(err)
			proof.Status = enums.PROOF_STATUS__FAILED
		} else {
			println(string(message))
			proof.Receipt = string(message)
			proof.Status = enums.PROOF_STATUS__SUCCEEDED
		}

		if err := proof.UpdateByProofID(d); err != nil {
			l.Error(err)
		}
	}()

	return &CreateProofRsp{
		Info: "proof is generating",
	}, nil
}
