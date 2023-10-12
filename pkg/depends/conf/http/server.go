package http

import (
	"context"
	"net/http"
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
	// Protocol support `http`, `https`, `unix`(http server based on unix domain socket), default `http`
	Protocol string `env:""`
	// Addr listen addr, default `0.0.0.0`. if based on a unix socket, set Addr as unix socket file abs path
	Addr string `env:""`
	// Port listen port
	Port int `env:",opt,expose"`
	// Spec document rel path
	Spec string `env:",opt,copy"`
	// HealthCheck path
	HealthCheck string `env:",opt,healthCheck"`
	// Debug if enable debug mode
	Debug *bool `env:""`

	ht       *httptransport.HttpTransport
	injector contextx.WithContext
	name     string
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
	if s.Protocol == "" {
		s.Protocol = "http"
	}
	if s.Addr == "" {
		s.Addr = "0.0.0.0"
	}
	if s.Port == 0 {
		switch s.Protocol {
		case "http":
			s.Port = 80
		case "https":
			s.Port = 443
		}
	}
	if s.Spec == "" {
		s.Spec = "./openapi.json"
	}
	if s.Debug == nil {
		s.Debug = ptrx.Ptr(true)
	}
	if s.HealthCheck == "" {
		s.HealthCheck = "http://:" + strconv.FormatInt(int64(s.Port), 10) + "/"
	}
	if s.ht == nil {
		modifiers := make([]httptransport.ServerModifier, 0)
		if s.Protocol == "http+unix" {
			modifiers = append(modifiers, func(srv *http.Server) error {
				srv.Addr = ":unix@" + s.Addr
				return nil
			})
		}

		s.ht = httptransport.NewHttpTransport(modifiers...)
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
