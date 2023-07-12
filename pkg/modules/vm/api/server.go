package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/conf/redis"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/modules/vm/api/async"
)

type Server struct {
	l   log.Logger
	cli *asynq.Client
	srv *asynq.Server
}

func (s *Server) Call(projectName string, data []byte) *http.Response {
	_, l := s.l.Start(context.Background(), "vm.api.Call")
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

func NewServer(redisConf *redis.Redis, mgrDB sqlx.DBExecutor, l log.Logger) *Server {
	router := newRouter()

	redisCli := asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%s:%d", redisConf.Host, redisConf.Port),
		Password: redisConf.Password.String(),
	}
	asyncCli := asynq.NewClient(redisCli)
	asyncSrv := asynq.NewServer(redisCli, asynq.Config{})
	mux := asynq.NewServeMux()

	mux.Handle(async.TaskNameApiCall, async.NewApiCallProcessor(router, asyncCli, l))
	mux.Handle(async.TaskNameApiResult, async.NewApiResultProcessor(mgrDB, l))

	_, l = l.Start(context.Background(), "vm.api.NewServer")
	defer l.End()

	if err := asyncSrv.Start(mux); err != nil {
		l.Fatal(err)
	}

	return &Server{
		l:   l,
		cli: asyncCli,
		srv: asyncSrv,
	}
}
