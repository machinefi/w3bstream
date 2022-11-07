package httptransport

import (
	"bytes"
	"context"
	"io"
	"mime"
	"net/http"
	"net/textproto"
	neturl "net/url"
	"reflect"
	"sort"

	"github.com/julienschmidt/httprouter"
	pkgerr "github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/transformer"
	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/validator"
	vldterr "github.com/machinefi/w3bstream/pkg/depends/kit/validator/errors"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/depends/x/reflectx"
)

type RequestTsfm struct {
	Type   reflect.Type
	Params map[string][]transformer.ReqParam
}

func (t *RequestTsfm) NewRequest(method, url string, v interface{}) (*http.Request, error) {
	return t.NewReqWithContext(context.Background(), method, url, v)
}

func (t *RequestTsfm) NewReqWithContext(ctx context.Context, method, url string, v interface{}) (*http.Request, error) {
	if v == nil {
		return http.NewRequestWithContext(ctx, method, url, nil)
	}

	typ := reflectx.DeRef(reflect.TypeOf(v))
	if t.Type != typ {
		return nil, pkgerr.Errorf(
			"unmatched request transformer, need %s but got %s", t.Type, typ,
		)
	}

	var (
		errs    = vldterr.NewErrorSet()
		params  = httprouter.Params{}
		query   = neturl.Values{}
		header  = http.Header{}
		cookies = neturl.Values{}
		body    = bytes.NewBuffer(nil)
	)

	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}
	rv = reflectx.Indirect(rv)

	for _, parameters := range t.Params {
		for i := range parameters {
			p := parameters[i]

			fv := p.FieldValue(rv)
			if !fv.IsValid() {
				continue
			}

			if p.In == "body" {
				if err := p.Tsf.EncodeTo(
					ctx,
					transformer.WriterWithHeader(body, header),
					fv,
				); err != nil {
					errs.AddErr(err, p.Name)
				}
				continue
			}

			writers := transformer.NewStringBuilders()
			if err := transformer.NewSuper(p.Tsf, &p.Option.CommonOption).
				EncodeTo(ctx, writers, fv); err != nil {
				errs.AddErr(err, p.Name)
				continue
			}

			values := writers.StringSlice()
			switch p.In {
			case "path":
				params = append(params, httprouter.Param{Key: p.Name, Value: values[0]})
			case "query":
				query[p.Name] = values
			case "header":
				header[textproto.CanonicalMIMEHeaderKey(p.Name)] = values
			case "cookie":
				cookies[p.Name] = values
			}
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	if len(params) > 0 {
		req = req.WithContext(
			contextx.WithValue(req.Context(), httprouter.ParamsKey, params),
		)
		req.URL.Path = transformer.NewPathnamePattern(req.URL.Path).Stringify(params)
	}

	if len(query) > 0 {
		if method == http.MethodGet && ShouldQueryInBodyForGet(ctx) {
			header.Set(
				httpx.HeaderContentType,
				mime.FormatMediaType(
					httpx.MIME_FORM_URLENCODED,
					map[string]string{"param": "value"},
				),
			)
			body = bytes.NewBufferString(query.Encode())
		} else {
			req.URL.RawQuery = query.Encode()
		}
	}

	req.Header = header

	if n := len(cookies); n > 0 {
		names := make([]string, n)
		i := 0
		for name := range cookies {
			names[i] = name
			i++
		}
		sort.Strings(names)

		for _, name := range names {
			values := cookies[name]
			for _, value := range values {
				req.AddCookie(&http.Cookie{Name: name, Value: value})
			}
		}
	}

	if n := int64(body.Len()); n != 0 {
		req.ContentLength = n
		rc := io.NopCloser(body)
		req.Body = rc
		req.GetBody = func() (io.ReadCloser, error) { return rc, nil }
	}

	return req, nil
}

func (t *RequestTsfm) DecodeAndValidate(ctx context.Context, ri httpx.RequestInfo, v interface{}) error {
	if err := t.DecodeFromRequestInfo(ctx, ri, v); err != nil {
		return err
	}
	return t.validate(v)
}

func (t *RequestTsfm) DecodeFromRequestInfo(ctx context.Context, ri httpx.RequestInfo, v interface{}) error {
	if with, ok := v.(httpx.WithFromRequestInfo); ok {
		if err := with.FromRequestInfo(ri); err != nil {
			if est := err.(interface {
				ToFieldErrors() statusx.ErrorFields
			}); ok {
				if errorFields := est.ToFieldErrors(); len(errorFields) > 0 {
					return (&badRequest{errorFields: errorFields}).Err()
				}
			}
			return err
		}
		return nil
	}

	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}

	if rv.Kind() != reflect.Ptr {
		return pkgerr.Errorf("decode target must be an ptr value")
	}

	rv = reflectx.Indirect(rv)

	if tpe := rv.Type(); tpe != t.Type {
		return pkgerr.Errorf(
			"unmatched request transformer, need %s but got %s",
			t.Type, tpe,
		)
	}

	errs := vldterr.NewErrorSet()

	for in := range t.Params {
		parameters := t.Params[in]

		for i := range parameters {
			param := parameters[i]

			if param.In == "body" {
				body := ri.Body()
				if err := param.Tsf.DecodeFrom(
					ctx,
					body,
					param.FieldValue(rv).Addr(),
					textproto.MIMEHeader(ri.Header()),
				); err != nil && err != io.EOF {
					errs.AddErr(err, vldterr.Location(param.In))
				}
				body.Close()
				continue
			}

			var values []string

			if param.In == "meta" {
				params := OperatorFactoryFromContext(ctx).Params
				if params != nil {
					values = params[param.Name]
				}
			} else {
				values = ri.Values(param.In, param.Name)
			}

			if len(values) > 0 {
				if err := transformer.NewSuper(
					param.Tsf,
					&param.Option.CommonOption,
				).DecodeFrom(
					ctx,
					transformer.NewStringReaders(values),
					param.FieldValue(rv).Addr(),
				); err != nil {
					errs.AddErr(err, vldterr.Location(param.In), param.Name)
				}
			}
		}
	}

	if errs.Err() == nil {
		return nil
	}

	return (&badRequest{errorFields: errs.ToErrorFields()}).Err()
}

func (t *RequestTsfm) validate(v interface{}) error {
	if self, ok := v.(validator.CanValidate); ok {
		if err := self.Validate(); err != nil {
			if est := err.(interface {
				ToFieldErrors() statusx.ErrorFields
			}); ok {
				if errorFields := est.ToFieldErrors(); len(errorFields) > 0 {
					return (&badRequest{errorFields: errorFields}).Err()
				}
			}
			return err
		}
		return nil
	}

	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}

	errSet := vldterr.NewErrorSet()

	for in := range t.Params {
		parameters := t.Params[in]

		for i := range parameters {
			param := parameters[i]

			if param.Validator != nil {
				if err := param.Validator.Validate(param.FieldValue(rv)); err != nil {
					if param.In == "body" {
						errSet.AddErr(err, vldterr.Location(param.In))
					} else {
						errSet.AddErr(err, vldterr.Location(param.In), param.Name)
					}
				}
			}
		}
	}

	br := &badRequest{errorFields: errSet.ToErrorFields()}

	// TODO deprecated
	if postValidator, ok := rv.Interface().(PostValidator); ok {
		postValidator.PostValidate(br)
	}

	if errSet.Err() == nil {
		return nil
	}

	return br.Err()
}

type PostValidator interface {
	PostValidate(badReqErr BadRequestError)
}
