package httptransport

import (
	"context"
	"net/http"
	"reflect"

	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/transformer"
	"github.com/machinefi/w3bstream/pkg/depends/kit/metax"
	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/depends/x/typesx"
)

type RouteHandler struct {
	*RequestTsfmFactory
	*HttpRouteMeta
	meta  *ServiceMeta
	tsfms []*RequestTsfm
}

func NewRouteHandler(sm *ServiceMeta, meta *HttpRouteMeta, rtf *RequestTsfmFactory) *RouteHandler {
	operatorFactories := meta.Metas

	if len(operatorFactories) == 0 {
		panic(errors.Errorf("missing valid operator"))
	}

	rts := make([]*RequestTsfm, len(operatorFactories))

	for i := range operatorFactories {
		opFactory := operatorFactories[i]
		rt, err := rtf.NewRequestTsfm(context.Background(), opFactory.Type)
		if err != nil {
			panic(err)
		}
		rts[i] = rt
	}

	return &RouteHandler{
		RequestTsfmFactory: rtf,
		HttpRouteMeta:      meta,
		meta:               sm,
		tsfms:              rts,
	}
}

// ServeHTTP impls http.Handler
func (hdl *RouteHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rid := hdl.Metas[len(hdl.Metas)-1].ID

	ctx := r.Context()
	ctx = ContextWithHttpRequest(ctx, r)
	ctx = ContextWithServiceMeta(ctx, *hdl.meta)
	ctx = ContextWithRouteMetaID(ctx, rid)

	// trace span, use operator as span name
	operator := []string{hdl.meta.String() + "/" + rid}
	ctx = metax.ContextWithMeta(ctx, metax.Meta{"operator": operator})

	rw.Header().Set("X-Meta", operator[0])

	ri := httpx.NewRequestInfo(r)

	for i := range hdl.Metas {
		m := hdl.Metas[i]
		if m.NoOutput {
			continue
		}
		op := m.New()

		ctx = ContextWithOperatorFactory(ctx, m.OperatorFactory)

		rt := hdl.tsfms[i]
		if rt != nil {
			err := rt.DecodeAndValidate(ctx, ri, op)
			if err != nil {
				hdl.writeErr(rw, r, err)
				return
			}
		}

		result, err := op.Output(ctx)

		if err != nil {
			hdl.writeErr(rw, r, err)
			return
		}

		if !m.IsLast {
			if c, ok := result.(context.Context); ok {
				ctx = c
			} else {
				// set result in context with key of operator name
				ctx = contextx.WithValue(ctx, m.ContextKey, result)
			}
			continue
		}
		hdl.write(rw, r, result)
	}
}

// resolve httpx.Encode (httpx.ResolveEncode)
func (hdl *RouteHandler) resolve(rsp *httpx.Response) (httpx.Encode, error) {
	tsfm, err := hdl.Tsfm.NewTransformer(
		context.Background(),
		typesx.FromReflectType(reflect.TypeOf(rsp.Value)),
		transformer.Option{MIME: rsp.ContentType},
	)
	if err != nil {
		return nil, err
	}
	return tsfm.EncodeTo, nil
}

// write response
func (hdl *RouteHandler) write(rw http.ResponseWriter, r *http.Request, rsp interface{}) {
	err := httpx.ResponseFrom(rsp).WriteTo(rw, r, hdl.resolve)
	if err != nil {
		hdl.writeErr(rw, r, err)
	}
}

func (hdl *RouteHandler) writeErr(rw http.ResponseWriter, r *http.Request, err error) {
	rsp, ok := err.(*httpx.Response)
	if !ok {
		rsp = httpx.ResponseFrom(err)
	}

	if se, ok := statusx.IsStatusErr(rsp.Unwrap()); ok {
		err := se.AppendSource(hdl.meta.String())

		if rwe, ok := rw.(ResponseWithError); ok {
			rwe.WriteError(err)
		}

		rsp.Value = err
	}

	errw := rsp.WriteTo(rw, r, hdl.resolve)
	if errw != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = rw.Write([]byte("write err failed:" + errw.Error()))
	}
}

type ResponseWithError interface {
	WriteError(err error)
}
