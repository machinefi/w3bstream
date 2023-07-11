package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/modules/vm/api/async"
	"github.com/machinefi/w3bstream/pkg/types"
)

type Server struct {
	ctx    context.Context
	router *gin.Engine
	cli    *asynq.Client
	srv    *asynq.Server
}

func (s *Server) Call(projectName string, data []byte) *http.Response {
	l := types.MustLoggerFromContext(s.ctx)
	_, l = l.Start(s.ctx, "vm.api.Call")
	defer l.End()

	task, err := async.NewApiCallTask(projectName, data)
	if err != nil {
		l.Error(errors.Wrap(err, "new api call task failed"))
		return &http.Response{
			StatusCode: http.StatusBadRequest,
		}
	}
	if _, err := s.cli.Enqueue(task); err != nil {
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

func newRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(ParamValidate())

	router.GET("/system/hello", hello)

	return router
}

func NewServer(ctx context.Context) *Server {
	router := newRouter()

	redisConf := types.MustRedisEndpointFromContext(ctx)
	redisCli := asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%s:%d", redisConf.Host, redisConf.Port),
		Password: redisConf.Password.String(),
	}
	asyncCli := asynq.NewClient(redisCli)
	asyncSrv := asynq.NewServer(redisCli, asynq.Config{})
	mux := asynq.NewServeMux()

	mux.Handle(async.TaskNameApiCall, async.NewApiCallProcessor(ctx, router, asyncCli))
	mux.Handle(async.TaskNameApiResult, async.NewApiResultProcessor(ctx))

	l := types.MustLoggerFromContext(ctx)
	_, l = l.Start(ctx, "vm.api.NewServer")
	defer l.End()

	if err := asyncSrv.Start(mux); err != nil {
		l.Fatal(err)
	}

	return &Server{
		ctx:    ctx,
		router: router,
		cli:    asyncCli,
		srv:    asyncSrv,
	}
}
