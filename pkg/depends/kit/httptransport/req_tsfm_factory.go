package httptransport

import (
	"context"
	"net/http"
	"reflect"
	"sync"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/transformer"
	"github.com/machinefi/w3bstream/pkg/depends/kit/validator"
	"github.com/machinefi/w3bstream/pkg/depends/x/reflectx"
	"github.com/machinefi/w3bstream/pkg/depends/x/typesx"
)

type RequestTsfmFactory struct {
	Vldt  validator.Factory
	Tsfm  transformer.Factory
	cache sync.Map
}

func NewRequestTsfmFactory(tsfm transformer.Factory, vldt validator.Factory) *RequestTsfmFactory {
	r := &RequestTsfmFactory{Vldt: vldt, Tsfm: tsfm}
	r.SetDefault()
	return r
}

func (f *RequestTsfmFactory) SetDefault() {
	if f.Vldt == nil {
		f.Vldt = validator.DefaultFactory
	}
	if f.Tsfm == nil {
		f.Tsfm = transformer.DefaultFactory
	}
}

func (f *RequestTsfmFactory) NewRequest(method, url string, v interface{}) (*http.Request, error) {
	return f.NewRequestWithContext(context.Background(), method, url, v)
}

func (f *RequestTsfmFactory) NewRequestWithContext(ctx context.Context, method, url string, v interface{}) (*http.Request, error) {
	if v == nil {
		return http.NewRequestWithContext(ctx, method, url, nil)
	}
	rt, err := f.NewRequestTsfm(AsRequestOut(ctx), reflect.TypeOf(v))
	if err != nil {
		return nil, err
	}
	return rt.NewReqWithContext(ctx, method, url, v)
}

func (f *RequestTsfmFactory) NewRequestTsfm(ctx context.Context, typ reflect.Type) (*RequestTsfm, error) {
	if v, ok := f.cache.Load(typ); ok {
		return v.(*RequestTsfm), nil
	}
	t, err := f.new(ctx, typ)
	if err != nil {
		return nil, err
	}
	f.cache.Store(typ, t)
	return t, nil
}

func (f *RequestTsfmFactory) new(ctx context.Context, typ reflect.Type) (*RequestTsfm, error) {
	rt := &RequestTsfm{}

	rt.Params = map[string][]transformer.ReqParam{}
	rt.Type = reflectx.DeRef(typ)

	ctx = transformer.ContextWithFactory(ctx, f.Tsfm)
	ctx = validator.ContextWithFactory(ctx, f.Vldt)

	err := transformer.EachReqParam(
		ctx, typesx.FromReflectType(rt.Type),
		func(rp *transformer.ReqParam) {
			if rp.In == "" {
				return
			}
			rt.Params[rp.In] = append(rt.Params[rp.In], *rp)
		},
	)

	return rt, err
}
