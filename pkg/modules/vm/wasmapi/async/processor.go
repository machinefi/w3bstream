package async

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/conf/log"
	confmq "github.com/machinefi/w3bstream/pkg/depends/conf/mq"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/event"
	apitypes "github.com/machinefi/w3bstream/pkg/modules/vm/wasmapi/types"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
	"github.com/machinefi/w3bstream/pkg/types/wasm/kvdb"
)

type ApiCallProcessor struct {
	l      log.Logger
	router *gin.Engine
	cli    *asynq.Client
}

func NewApiCallProcessor(l log.Logger, router *gin.Engine, cli *asynq.Client) *ApiCallProcessor {
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
	req = req.WithContext(contextx.WithContextCompose(
		types.WithProjectContext(payload.Project),
		wasm.WithChainClientContext(payload.ChainClient),
		types.WithLoggerContext(p.l),
	)(context.Background()))

	respRecorder := httptest.NewRecorder()
	p.router.ServeHTTP(respRecorder, req)

	_, l := p.l.Start(ctx, "wasmapi.ProcessTaskApiCall")
	defer l.End()
	prjName := payload.Project.ProjectName.Name
	l = l.WithValues("ProjectName", prjName)

	apiResp, err := ConvHttpResponse(req.Header, respRecorder.Result())
	if err != nil {
		l.Error(errors.Wrap(err, "conv http response failed"))
		return fmt.Errorf("conv http response failed: %v: %w", err, asynq.SkipRetry)
	}

	// no content need return to caller
	if apiResp.StatusCode == http.StatusNoContent {
		return nil
	}

	apiRespJson, err := json.Marshal(apiResp)
	if err != nil {
		l.Error(errors.Wrap(err, "encode http response failed"))
		return fmt.Errorf("encode http response failed: %v: %w", err, asynq.SkipRetry)
	}

	eventType := req.Header.Get("eventType")

	task, err := newApiResultTask(prjName, eventType, apiRespJson)
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
	kv    *kvdb.RedisDB
	tasks *confmq.Config
}

func NewApiResultProcessor(l log.Logger, mgrDB sqlx.DBExecutor, kv *kvdb.RedisDB, tasks *confmq.Config) *ApiResultProcessor {
	return &ApiResultProcessor{
		l:     l,
		kv:    kv,
		mgrDB: mgrDB,
		tasks: tasks,
	}
}

func (p *ApiResultProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	payload := apiResultPayload{}
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	ctx = contextx.WithContextCompose(
		confmq.WithMqContext(p.tasks),
		types.WithLoggerContext(p.l),
		types.WithMgrDBExecutorContext(p.mgrDB),
		kvdb.WithRedisDBKeyContext(p.kv),
		types.WithProjectContext(&models.Project{
			ProjectName: models.ProjectName{Name: payload.ProjectName}},
		),
	)(ctx)

	_, l := p.l.Start(ctx, "wasmapi.ProcessTaskApiResult")
	defer l.End()

	if _, err := event.HandleEvent(ctx, payload.EventType, payload.Data); err != nil {
		l.Error(errors.Wrap(err, "send event failed"))
		return err
	}

	return nil
}

func ConvHttpResponse(reqHeader http.Header, resp *http.Response) (*apitypes.HttpResponse, error) {
	var body []byte
	var err error
	if resp.Body != nil {
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
	}

	respHeader := resp.Header
	for k, v := range reqHeader {
		if k == "Content-Type" {
			continue
		}
		respHeader[k] = v
	}

	return &apitypes.HttpResponse{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Proto:      resp.Proto,
		Header:     respHeader,
		Body:       body,
	}, nil
}
