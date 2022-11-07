package httptransport

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/handlers"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/transformer"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/depends/kit/validator"
)

func MiddlewareChain(mw ...HttpMiddleware) HttpMiddleware {
	return func(final http.Handler) http.Handler {
		last := final
		for i := len(mw) - 1; i >= 0; i-- {
			last = mw[i](last)
		}
		return last
	}
}

type HttpMiddleware func(http.Handler) http.Handler

func NewHttpTransport(modifiers ...ServerModifier) *HttpTransport {
	return &HttpTransport{Modifiers: modifiers}
}

type HttpTransport struct {
	ServiceMeta
	Port        int
	Modifiers   []ServerModifier    // for modifying http.Server
	Middlewares []HttpMiddleware    // Middlewares https://github.com/gorilla/handlers
	Vldt        validator.Factory   // Vldt validator factory
	Tsfm        transformer.Factory // transformer mgr for parameter transforming
	CertFile    string
	KeyFile     string
	httpRouter  *httprouter.Router
}

type ServerModifier func(server *http.Server) error

func (t *HttpTransport) SetDefault() {
	t.ServiceMeta.SetDefault()

	if t.Vldt == nil {
		t.Vldt = validator.DefaultFactory
	}

	if t.Tsfm == nil {
		t.Tsfm = transformer.DefaultFactory
	}

	if t.Middlewares == nil {
		t.Middlewares = []HttpMiddleware{handlers.LogHandler()}
	}

	if t.Port == 0 {
		t.Port = 80
	}
}

func (t *HttpTransport) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	t.httpRouter.ServeHTTP(w, req)
}

func (t *HttpTransport) Serve(router *kit.Router) error {
	return t.ServeContext(context.Background(), router)
}

func (t *HttpTransport) ServeContext(ctx context.Context, router *kit.Router) error {
	t.SetDefault()

	logger := log.FromContext(ctx)

	t.httpRouter = t.toHttpRouter(router)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", t.Port),
		Handler: MiddlewareChain(t.Middlewares...)(t),
	}

	for i := range t.Modifiers {
		if err := t.Modifiers[i](srv); err != nil {
			logger.Fatal(err)
		}
	}

	go func() {
		outputln("%s listen on %s", t.ServiceMeta, srv.Addr)

		if t.CertFile != "" && t.KeyFile != "" {
			if err := srv.ListenAndServeTLS(t.CertFile, t.KeyFile); err != nil {
				if err == http.ErrServerClosed {
					logger.Error(err)
				} else {
					logger.Fatal(err)
				}
			}
			return
		}

		if err := srv.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				logger.Error(err)
			} else {
				logger.Fatal(err)
			}
		}
	}()

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)
	<-stopCh

	timeout := 10 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	logger.Info("shutdowning in %s", timeout)

	return srv.Shutdown(ctx)
}

func (t *HttpTransport) toHttpRouter(rt *kit.Router) *httprouter.Router {
	routes := rt.Routes()

	if len(routes) == 0 {
		panic(errors.Errorf(
			"need to register Operator to Router %#v before serve", rt,
		))
	}

	metas := make([]*HttpRouteMeta, len(routes))
	for i := range routes {
		metas[i] = NewHttpRouteMeta(routes[i])
	}

	router := httprouter.New()

	sort.Slice(metas, func(i, j int) bool {
		return metas[i].Key() < metas[j].Key()
	})

	for i := range metas {
		route := metas[i]
		route.Log()

		if err := tryCatch(func() {
			router.HandlerFunc(
				route.Method(),
				route.Path(),
				NewRouteHandler(
					&t.ServiceMeta,
					route,
					NewRequestTsfmFactory(t.Tsfm, t.Vldt),
				).ServeHTTP,
			)
		}); err != nil {
			panic(errors.Errorf("register http route `%s` failed: %s", route, err))
		}
	}

	return router
}
