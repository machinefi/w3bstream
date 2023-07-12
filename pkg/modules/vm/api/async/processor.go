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

	"github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/event"
	"github.com/machinefi/w3bstream/pkg/types"
)

type ApiCallProcessor struct {
	l      log.Logger
	router *gin.Engine
	cli    *asynq.Client
}

func NewApiCallProcessor(router *gin.Engine, cli *asynq.Client, l log.Logger) *ApiCallProcessor {
	return &ApiCallProcessor{
		l:      l,
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

	_, l := p.l.Start(context.Background(), "vm.api.ProcessTaskApiCall")
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
	l     log.Logger
	mgrDB sqlx.DBExecutor
}

func NewApiResultProcessor(mgrDB sqlx.DBExecutor, l log.Logger) *ApiResultProcessor {
	return &ApiResultProcessor{
		l:     l,
		mgrDB: mgrDB,
	}
}

func (p *ApiResultProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	payload := apiResultPayload{}
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	ctx = types.WithMgrDBExecutor(ctx, p.mgrDB)
	ctx = types.WithLogger(ctx, p.l)
	ctx = types.WithProject(ctx, &models.Project{
		ProjectName: models.ProjectName{Name: payload.ProjectName}},
	)

	_, l := p.l.Start(ctx, "vm.api.ProcessTaskApiResult")
	defer l.End()

	if _, err := event.HandleEvent(ctx, payload.EventType, payload.Data); err != nil {
		l.Error(errors.Wrap(err, "send event failed"))
		return err
	}

	return nil
}
