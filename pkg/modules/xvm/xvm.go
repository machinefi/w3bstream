package xvm

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
)

func CreateRisc0Vm(ctx context.Context, req *CreateRisc0VmReq) (*CreateRisc0VmRsp, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fields := map[string]string{
		"image_id": req.ImageId,
		"elf":      req.Elf,
	}
	for key, val := range fields {
		_ = w.WriteField(key, val)
	}

	params := req.Params
	for _, param := range params {
		_ = w.WriteField("params", param)
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
	request, err := http.NewRequest("POST", "http://127.0.0.1:3000/api/prove_file", &b)
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
