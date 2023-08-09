package http

import (
	"context"
	"strconv"

	"go.opentelemetry.io/otel"

	"github.com/machinefi/w3bstream/pkg/depends/conf/http/mws"
	"github.com/machinefi/w3bstream/pkg/depends/conf/logger"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/depends/x/ptrx"
)

var middlewares []httptransport.HttpMiddleware

// WithMiddlewares for custom
func WithMiddlewares(ms ...httptransport.HttpMiddleware) {
	middlewares = append(middlewares, ms...)
}

type Server struct {
	Port        int                          `env:",opt,expose"`
	Spec        string                       `env:",opt,copy"`
	HealthCheck string                       `env:",opt,healthCheck"`
	Debug       *bool                        `env:""`
	ht          *httptransport.HttpTransport `env:"-"`
	injector    contextx.WithContext         `env:"-"`
	name        string
}

func (s Server) WithContextInjector(injector contextx.WithContext) *Server {
	s.injector = injector
	return &s
}

func (s Server) WithName(name string) *Server {
	s.name = name
	return &s
}

func (s *Server) LivenessCheck() map[string]string {
	statuses := map[string]string{}

	if s.ht != nil {
		statuses[s.ht.ServiceMeta.String()] = "ok"
	}

	return statuses
}

func (s *Server) SetDefault() {
	if s.Port == 0 {
		s.Port = 80
	}

	if s.Spec == "" {
		s.Spec = "./openapi.json"
	}

	if s.Debug == nil {
		s.Debug = ptrx.Bool(true)
	}

	if s.HealthCheck == "" {
		s.HealthCheck = "http://:" + strconv.FormatInt(int64(s.Port), 10) + "/"
	}

	if s.ht == nil {
		s.ht = httptransport.NewHttpTransport()
		s.ht.SetDefault()
	}
}

func (s *Server) Serve(router *kit.Router) error {
	if s.ht == nil {
		s.ht = httptransport.NewHttpTransport()
		s.ht.SetDefault()
	}

	tr := otel.Tracer(s.name)
	ht := s.ht
	ht.Port = s.Port

	ht.Middlewares = []httptransport.HttpMiddleware{}
	ht.Middlewares = append(ht.Middlewares, middlewares...)
	ht.Middlewares = append(ht.Middlewares,
		mws.DefaultCORS(),
		mws.HealthCheckHandler(),
		mws.MetricsHandler(),
		TraceLogHandler(tr),
		NewContextInjectorMw(s.injector),
	)
	if s.Debug != nil && *s.Debug {
		ht.Middlewares = append(ht.Middlewares, mws.PProfHandler(*s.Debug))
	}

	ctx, _ := logger.NewSpanContext(context.Background(), s.name)

	return s.ht.ServeContext(ctx, router)
}

func (s *Server) Shutdown() {
	_ = s.ht.Shutdown(context.Background())
}
