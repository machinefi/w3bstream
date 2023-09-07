package wasmapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/pkg/errors"

	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/conf/redis"
	"github.com/machinefi/w3bstream/pkg/depends/kit/mq"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	optypes "github.com/machinefi/w3bstream/pkg/modules/operator/pool/types"
	"github.com/machinefi/w3bstream/pkg/modules/vm/wasmapi/async"
	"github.com/machinefi/w3bstream/pkg/modules/vm/wasmapi/handler"
	apitypes "github.com/machinefi/w3bstream/pkg/modules/vm/wasmapi/types"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm/kvdb"
)

type Server struct {
	router *gin.Engine
	cli    *asynq.Client
	srv    *asynq.Server
}

func (s *Server) Call(ctx context.Context, data []byte) *apitypes.HttpResponse {
	l := types.MustLoggerFromContext(ctx)
	_, l = l.Start(ctx, "wasmapi.Call")
	defer l.End()

	apiReq := apitypes.HttpRequest{}
	if err := json.Unmarshal(data, &apiReq); err != nil {
		l.Error(errors.Wrap(err, "http request illegal format"))
		return &apitypes.HttpResponse{
			StatusCode: http.StatusBadRequest,
		}
	}
	req, err := http.NewRequestWithContext(ctx, apiReq.Method, apiReq.Url, bytes.NewReader(apiReq.Body))
	if err != nil {
		l.Error(errors.Wrap(err, "build http request failed"))
		return &apitypes.HttpResponse{
			StatusCode: http.StatusBadRequest,
		}
	}
	req.Header = apiReq.Header

	respRecorder := httptest.NewRecorder()
	s.router.ServeHTTP(respRecorder, req)

	resp, err := async.ConvHttpResponse(apiReq.Header, respRecorder.Result())
	if err != nil {
		l.Error(errors.Wrap(err, "conv http response failed"))
		return &apitypes.HttpResponse{
			StatusCode: http.StatusInternalServerError,
		}
	}
	return resp
}

func (s *Server) Shutdown() {
	s.srv.Shutdown()
}

func newRouter(mgrDB sqlx.DBExecutor, chainConf *types.ChainConfig, opPool optypes.Pool, sfid confid.SFIDGenerator, asyncCli *asynq.Client) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(handler.ParamValidate())

	handlers := handler.New(mgrDB, chainConf, opPool, sfid, asyncCli)

	router.GET("/system/hello", handlers.Hello)
	router.GET("/system/hello/async", handlers.HelloAsync)
	router.GET("/system/read_tx", handlers.ReadTx)
	router.GET("/system/read_tx/async", handlers.ReadTxAsync)
	router.POST("/system/send_tx", handlers.SendTx)
	router.POST("/system/send_tx/async", handlers.SendTxAsync)
	router.POST("/system/send_tx/async/state", handlers.SendTxAsyncStateCheck)

	return router
}

func NewServer(l log.Logger, redisConf *redis.Redis, mgrDB sqlx.DBExecutor, kv *kvdb.RedisDB, chainConf *types.ChainConfig, tb *mq.TaskBoard, tw *mq.TaskWorker, opPool optypes.Pool, sfid confid.SFIDGenerator) (*Server, error) {

	redisCli := asynq.RedisClientOpt{
		Network:      redisConf.Protocol,
		Addr:         fmt.Sprintf("%s:%d", redisConf.Host, redisConf.Port),
		Password:     redisConf.Password.String(),
		ReadTimeout:  time.Duration(redisConf.ReadTimeout),
		WriteTimeout: time.Duration(redisConf.WriteTimeout),
		DialTimeout:  time.Duration(redisConf.ConnectTimeout),
		DB:           redisConf.DB,
	}
	asyncCli := asynq.NewClient(redisCli)
	asyncSrv := asynq.NewServer(redisCli, asynq.Config{})

	router := newRouter(mgrDB, chainConf, opPool, sfid, asyncCli)

	mux := asynq.NewServeMux()
	mux.Handle(async.TaskNameApiCall, async.NewApiCallProcessor(l, router, asyncCli))
	mux.Handle(async.TaskNameApiResult, async.NewApiResultProcessor(l, mgrDB, kv, tb, tw))

	if err := asyncSrv.Start(mux); err != nil {
		return nil, err
	}

	return &Server{
		router: router,
		cli:    asyncCli,
		srv:    asyncSrv,
	}, nil
}
