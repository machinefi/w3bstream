package transformer

import (
	"context"
	"io"
	"net/textproto"
	"net/url"
	"reflect"

	pkgerr "github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	vldterr "github.com/machinefi/w3bstream/pkg/depends/kit/validator/errors"
	"github.com/machinefi/w3bstream/pkg/depends/x/reflectx"
	"github.com/machinefi/w3bstream/pkg/depends/x/typesx"
)

func init() { DefaultFactory.Register(&URLEncoded{}) }

/*
URLEncoded for application/x-www-form-urlencoded

	var s = struct {
		Username string `name:"username"`
		Nickname string `name:"username,omitempty"`
		Tags []string   `name:"tag"`
	}{
		Username: "name",
		Tags: []string{"1","2"},
	}

will transform to

	username=name&tag=1&tag=2
*/
type URLEncoded struct{ *FlattenParams }

func (URLEncoded) Names() []string {
	return []string{httpx.MIME_FORM_URLENCODED, "form", "urlencoded", "url-encoded"}
}

func (t URLEncoded) String() string { return httpx.MIME_FORM_URLENCODED }

func (URLEncoded) NamedByTag() string { return "name" }

func (URLEncoded) New(ctx context.Context, typ typesx.Type) (Transformer, error) {
	tsf := &URLEncoded{}

	typ = typesx.DeRef(typ)
	if typ.Kind() != reflect.Struct {
		return nil, pkgerr.Errorf(
			"content transformer `%s` should be used for struct type",
			tsf,
		)
	}

	tsf.FlattenParams = &FlattenParams{}

	if err := tsf.FlattenParams.CollectParams(ctx, typ); err != nil {
		return nil, err
	}

	return tsf, nil
}

func (t *URLEncoded) EncodeTo(ctx context.Context, w io.Writer, v interface{}) error {
	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}
	rv = reflectx.Indirect(rv)

	values := url.Values{}
	errs := vldterr.NewErrorSet()

	for i := range t.Params {
		p := t.Params[i]

		if p.Tsf != nil {
			field := p.FieldValue(rv)
			builders := NewStringBuilders()
			if err := NewSuper(p.Tsf, &p.Option.CommonOption).
				EncodeTo(ctx, builders, field); err != nil {
				errs.AddErr(err, p.Name)
				continue
			}
			values[p.Name] = builders.StringSlice()
		}
	}

	if err := errs.Err(); err != nil {
		return err
	}

	httpx.MaybeWriteHeader(
		ctx, w, t.Names()[0],
		map[string]string{
			"param": "value",
		},
	)
	_, err := w.Write([]byte(values.Encode()))
	return err
}

func (t *URLEncoded) DecodeFrom(ctx context.Context, r io.Reader, v interface{}, headers ...textproto.MIMEHeader) error {
	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}

	if rv.Kind() != reflect.Ptr {
		return pkgerr.New("decode target must be ptr value")
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	values, err := url.ParseQuery(string(data))
	if err != nil {
		return err
	}

	errs := vldterr.NewErrorSet()
	for i := range t.Params {
		p := t.Params[i]
		fvs := values[p.Name]

		if p.Tsf == nil || len(fvs) == 0 {
			continue
		}
		if err := NewSuper(p.Tsf, &p.Option.CommonOption).
			DecodeFrom(ctx, NewStringReaders(fvs), p.FieldValue(rv).Addr()); err != nil {
			errs.AddErr(err, p.Name)
			continue
		}
	}
	return errs.Err()
}
