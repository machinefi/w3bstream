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

	"github.com/machinefi/w3bstream/pkg/depends/conf/logger"
	confmq "github.com/machinefi/w3bstream/pkg/depends/conf/mq"
	confredis "github.com/machinefi/w3bstream/pkg/depends/conf/redis"
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
	router *gin.Engine
	cli    *asynq.Client
}

func NewApiCallProcessor(router *gin.Engine, cli *asynq.Client) *ApiCallProcessor {
	return &ApiCallProcessor{
		router: router,
		cli:    cli,
	}
}

func (p *ApiCallProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	ctx, l := logger.NewSpanContext(ctx, "vm.ApiCall.ProcessTask")
	defer l.End()

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
	)(context.Background()))

	respRecorder := httptest.NewRecorder()
	p.router.ServeHTTP(respRecorder, req)

	prjName := payload.Project.ProjectName.Name
	l = l.WithValues("prj", prjName)

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
	mgrDB sqlx.DBExecutor
	kv    *kvdb.RedisDB
	redis *confredis.Redis
	tasks *confmq.Config
}

func NewApiResultProcessor(mgrDB sqlx.DBExecutor, kv *kvdb.RedisDB, tasks *confmq.Config, redis *confredis.Redis) *ApiResultProcessor {
	return &ApiResultProcessor{
		kv:    kv,
		mgrDB: mgrDB,
		tasks: tasks,
		redis: redis,
	}
}

func (p *ApiResultProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	ctx, l := logger.NewSpanContext(ctx, "vm.ApiResult.ProcessTask")
	defer l.End()

	payload := apiResultPayload{}
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	ctx = contextx.WithContextCompose(
		confmq.WithMqContext(p.tasks),
		types.WithMgrDBExecutorContext(p.mgrDB),
		kvdb.WithRedisDBKeyContext(p.kv),
		types.WithProjectContext(&models.Project{
			ProjectName: models.ProjectName{Name: payload.ProjectName}},
		),
		types.WithRedisEndpointContext(p.redis),
	)(ctx)

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
