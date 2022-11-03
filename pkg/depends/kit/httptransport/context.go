package httptransport

import (
	"context"
	"net/http"
	"os"

	"github.com/machinefi/w3bstream/pkg/depends/base/consts"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
)

// TODO move to conf

type ServiceMeta struct {
	Name    string
	Version string
}

func (s *ServiceMeta) SetDefault() {
	if s.Name == "" {
		s.Name = os.Getenv(consts.EnvProjectName)
	}
	if s.Version == "" {
		s.Version = os.Getenv(consts.EnvProjectVersion)
	}
}

func (s ServiceMeta) String() string {
	if s.Version == "" {
		return s.Name
	}
	return s.Name + "@" + s.Version
}

type (
	keyHttpRequest     struct{} // keyHttpRequest pass original *http.Request
	keyServiceMeta     struct{} // keyServiceMeta pass service meta
	keyRouteMetaID     struct{}
	ckOperatorFactory  struct{}
	ckQueryInBoyForGet struct{}
	ckForRequestOut    struct{}
)

func ContextWithHttpRequest(ctx context.Context, req *http.Request) context.Context {
	return contextx.WithValue(ctx, keyHttpRequest{}, req)
}

func HttpRequestFromContext(ctx context.Context) *http.Request {
	r, _ := ctx.Value(keyHttpRequest{}).(*http.Request)
	return r
}

func ContextWithServiceMeta(ctx context.Context, meta ServiceMeta) context.Context {
	return contextx.WithValue(ctx, keyServiceMeta{}, meta)
}

func ServiceMetaFromContext(ctx context.Context) ServiceMeta {
	r, _ := ctx.Value(keyServiceMeta{}).(ServiceMeta)
	return r
}

func ContextWithRouteMetaID(ctx context.Context, id string) context.Context {
	return contextx.WithValue(ctx, keyRouteMetaID{}, id)
}

func OperationIDFromContext(ctx context.Context) string {
	r, _ := ctx.Value(keyRouteMetaID{}).(string)
	return r
}

func ContextWithOperatorFactory(ctx context.Context, factory *kit.OperatorFactory) context.Context {
	return contextx.WithValue(ctx, ckOperatorFactory{}, factory)
}

func OperatorFactoryFromContext(ctx context.Context) *kit.OperatorFactory {
	r, _ := ctx.Value(ckOperatorFactory{}).(*kit.OperatorFactory)
	return r
}

func EnableQueryInBodyForGet(ctx context.Context) context.Context {
	return contextx.WithValue(ctx, ckQueryInBoyForGet{}, true)
}

func ShouldQueryInBodyForGet(ctx context.Context) bool {
	v, ok := ctx.Value(ckQueryInBoyForGet{}).(bool)
	return ok && v
}

func AsRequestOut(ctx context.Context) context.Context {
	return contextx.WithValue(ctx, ckForRequestOut{}, true)
}

func IsRequestOut(ctx context.Context) bool {
	if ctx == nil {
		return false
	}
	if t, ok := ctx.Value(ckForRequestOut{}).(bool); ok {
		return t
	}
	return false
}
