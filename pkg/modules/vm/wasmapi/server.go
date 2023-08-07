package wasmapi

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/conf/redis"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/modules/vm/wasmapi/async"
	"github.com/machinefi/w3bstream/pkg/modules/vm/wasmapi/handler"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
	"github.com/machinefi/w3bstream/pkg/types/wasm/kvdb"
)

type Server struct {
	cli *asynq.Client
	srv *asynq.Server
}

func (s *Server) Call(ctx context.Context, data []byte) *http.Response {
	l := types.MustLoggerFromContext(ctx)
	_, l = l.Start(ctx, "wasmapi.Call")
	defer l.End()

	prj := types.MustProjectFromContext(ctx)
	chainCli := wasm.MustChainClientFromContext(ctx)
	task, err := async.NewApiCallTask(prj, chainCli, data)
	if err != nil {
		l.Error(errors.Wrap(err, "new api call task failed"))
		return &http.Response{
			StatusCode: http.StatusBadRequest,
		}
	}
	if _, err := s.cli.EnqueueContext(ctx, task); err != nil {
		l.Error(errors.Wrap(err, "could not enqueue task"))
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
		}
	}

	return &http.Response{
		StatusCode: http.StatusOK,
	}
}

func (s *Server) Shutdown() {
	s.srv.Shutdown()
}

func newRouter(mgrDB sqlx.DBExecutor, chainConf *types.ChainConfig) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(handler.ParamValidate())

	handlers := handler.New(mgrDB, chainConf)

	router.GET("/system/hello", handlers.Hello)
	router.GET("/system/read_tx", handlers.ReadTx)
	router.POST("/system/send_tx", handlers.SendTx)

	return router
}

func NewServer(l log.Logger, redisConf *redis.Redis, mgrDB sqlx.DBExecutor, kv *kvdb.RedisDB, chainConf *types.ChainConfig) (*Server, error) {
	router := newRouter(mgrDB, chainConf)

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
	mux := asynq.NewServeMux()

	mux.Handle(async.TaskNameApiCall, async.NewApiCallProcessor(l, router, asyncCli))
	mux.Handle(async.TaskNameApiResult, async.NewApiResultProcessor(l, mgrDB, kv))

	if err := asyncSrv.Start(mux); err != nil {
		return nil, err
	}

	return &Server{
		cli: asyncCli,
		srv: asyncSrv,
	}, nil
}
