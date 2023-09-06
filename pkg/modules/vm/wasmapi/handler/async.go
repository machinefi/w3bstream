package handler

import (
	"bytes"
	"net/http"
	"path"

	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/modules/vm/wasmapi/async"
)

func (h *Handler) setAsync(req *http.Request) error {
	req.URL.Path = path.Join(req.URL.Path, "async")

	var buf bytes.Buffer
	if err := req.Write(&buf); err != nil {
		return errors.Wrap(err, "http request write to buffer failed")
	}

	task, err := async.NewApiCallTask(buf.Bytes())
	if err != nil {
		return errors.Wrap(err, "new api call task failed")
	}
	if _, err := h.asyncCli.Enqueue(task); err != nil {
		return errors.Wrap(err, "could not enqueue task")
	}
	return nil
}
