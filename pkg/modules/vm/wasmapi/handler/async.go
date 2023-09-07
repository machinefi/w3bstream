package handler

import (
	"bytes"
	"io"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/modules/vm/wasmapi/async"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

func (h *Handler) setAsync(c *gin.Context) error {
	if cb, ok := c.Get(gin.BodyBytesKey); ok {
		if cbb, ok := cb.([]byte); ok {
			c.Request.Body = io.NopCloser(bytes.NewReader(cbb))
		}
	}
	req := c.Request
	req.URL.Path = path.Join(req.URL.Path, "async")

	var buf bytes.Buffer
	if err := req.Write(&buf); err != nil {
		return errors.Wrap(err, "http request write to buffer failed")
	}

	prj := types.MustProjectFromContext(req.Context())
	chainCli := wasm.MustChainClientFromContext(req.Context())

	task, err := async.NewApiCallTask(prj, chainCli, buf.Bytes())
	if err != nil {
		return errors.Wrap(err, "new api call task failed")
	}
	if _, err := h.asyncCli.Enqueue(task); err != nil {
		return errors.Wrap(err, "could not enqueue task")
	}
	return nil
}

func (h *Handler) setAsyncAdvance(c *gin.Context, path string, after time.Duration) error {
	if cb, ok := c.Get(gin.BodyBytesKey); ok {
		if cbb, ok := cb.([]byte); ok {
			c.Request.Body = io.NopCloser(bytes.NewReader(cbb))
		}
	}
	req := c.Request
	req.URL.Path = path

	var buf bytes.Buffer
	if err := req.Write(&buf); err != nil {
		return errors.Wrap(err, "http request write to buffer failed")
	}

	prj := types.MustProjectFromContext(req.Context())
	chainCli := wasm.MustChainClientFromContext(req.Context())

	task, err := async.NewApiCallTask(prj, chainCli, buf.Bytes())
	if err != nil {
		return errors.Wrap(err, "new api call task failed")
	}
	if _, err := h.asyncCli.Enqueue(task, asynq.ProcessIn(after)); err != nil {
		return errors.Wrap(err, "could not enqueue task")
	}
	return nil
}
