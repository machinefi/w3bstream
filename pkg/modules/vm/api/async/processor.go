package async

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/event"
	"github.com/machinefi/w3bstream/pkg/types"
)

type ApiCallProcessor struct {
	sysCtx context.Context
	router *gin.Engine
	cli    *asynq.Client
}

func NewApiCallProcessor(ctx context.Context, router *gin.Engine, cli *asynq.Client) *ApiCallProcessor {
	return &ApiCallProcessor{
		sysCtx: ctx,
		router: router,
		cli:    cli,
	}
}

func (p *ApiCallProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	payload := apiCallPayload{}
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	req, err := http.ReadRequest(bufio.NewReader(bytes.NewBuffer(payload.Data)))
	if err != nil {
		return fmt.Errorf("http.ReadRequest failed: %v: %w", err, asynq.SkipRetry)
	}

	resp := httptest.NewRecorder()
	p.router.ServeHTTP(resp, req)

	l := types.MustLoggerFromContext(p.sysCtx)
	_, l = l.Start(p.sysCtx, "vm.api.ProcessTaskApiCall")
	defer l.End()
	l = l.WithValues("ProjectName", payload.ProjectName)

	wbuf := bytes.Buffer{}
	if err := resp.Result().Write(&wbuf); err != nil {
		l.Error(errors.Wrap(err, "encode http response failed"))
		return fmt.Errorf("encode http response failed: %v: %w", err, asynq.SkipRetry)
	}

	eventType := req.Header.Get("eventType")
	if eventType == "" {
		l.Error(errors.New("miss eventType"))
		return fmt.Errorf("miss eventType, projectName %v: %w", payload.ProjectName, asynq.SkipRetry)
	}

	task, err := newApiResultTask(payload.ProjectName, eventType, wbuf.Bytes())
	if err != nil {
		l.Error(errors.Wrap(err, "new api result task failed"))
		return fmt.Errorf("new api result task failed: %v: %w", err, asynq.SkipRetry)
	}
	if _, err := p.cli.Enqueue(task); err != nil {
		l.Error(errors.Wrap(err, "could not enqueue task"))
		return fmt.Errorf("could not enqueue task: %v: %w", err, asynq.SkipRetry)
	}

	return nil
}

type ApiResultProcessor struct {
	sysCtx context.Context
}

func NewApiResultProcessor(ctx context.Context) *ApiResultProcessor {
	return &ApiResultProcessor{
		sysCtx: ctx,
	}
}

func (p *ApiResultProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	payload := apiResultPayload{}
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	sysCtx := types.WithProject(p.sysCtx, &models.Project{
		ProjectName: models.ProjectName{Name: payload.ProjectName}},
	)

	l := types.MustLoggerFromContext(sysCtx)
	_, l = l.Start(p.sysCtx, "vm.api.ProcessTaskApiResult")
	defer l.End()

	if _, err := event.HandleEvent(sysCtx, payload.EventType, payload.Data); err != nil {
		l.Error(errors.Wrap(err, "send event failed"))
		return err
	}

	return nil
}
