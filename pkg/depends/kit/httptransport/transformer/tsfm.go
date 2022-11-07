package transformer

import (
	"context"
	"io"
	"net/textproto"
	"net/url"
	"reflect"
	"sort"

	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/depends/x/mapx"
	"github.com/machinefi/w3bstream/pkg/depends/x/typesx"
)

type Transformer interface {
	Names() []string
	New(context.Context, typesx.Type) (Transformer, error)
	EncodeTo(context.Context, io.Writer, any) error
	DecodeFrom(context.Context, io.Reader, any, ...textproto.MIMEHeader) error
}

type Option struct {
	Name string
	MIME string
	CommonOption
}

func (o Option) String() string {
	values := url.Values{}
	if o.Name != "" {
		values.Add("Name", o.Name)
	}
	if o.MIME != "" {
		values.Add("MIME", o.MIME)
	}
	if o.Omitempty {
		values.Add("Omitempty", "true")
	}
	if o.Explode {
		values.Add("Explode", "true")
	}
	return values.Encode()
}

type CommonOption struct {
	Omitempty bool // Omitempty should be ignored when value is empty
	Explode   bool // Explode for raw uint8/byte slice/array
}

type TsfmAndOption struct {
	Transformer
	Option Option
}

type Factory interface {
	NewTransformer(context.Context, typesx.Type, Option) (Transformer, error)
}

type ckFactory struct{}

func ContextWithFactory(ctx context.Context, f Factory) context.Context {
	return contextx.WithValue(ctx, ckFactory{}, f)
}

func FactoryFromContext(ctx context.Context) Factory {
	if f, ok := ctx.Value(ckFactory{}).(Factory); ok {
		return f
	}
	return DefaultFactory
}

func NewTransformer(ctx context.Context, t typesx.Type, opt Option) (Transformer, error) {
	return FactoryFromContext(ctx).NewTransformer(ctx, t, opt)
}

type factory struct {
	set   map[string]Transformer
	cache *mapx.Map[string, Transformer]
}

var DefaultFactory = &factory{}

func NewFactory() *factory {
	return &factory{
		set:   make(map[string]Transformer),
		cache: mapx.New[string, Transformer](),
	}
}

func (f *factory) Register(tsfms ...Transformer) {
	if f.set == nil {
		f.set = map[string]Transformer{}
	}
	if f.cache == nil {
		f.cache = mapx.New[string, Transformer]()
	}
	for _, t := range tsfms {
		for _, name := range t.Names() {
			f.set[name] = t
		}
	}
}

func (f *factory) NewTransformer(ctx context.Context, t typesx.Type, opt Option) (Transformer, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	key := typesx.FullTypeName(t) + opt.String()

	if v, ok := f.cache.Load(key); ok {
		return v.(Transformer), nil
	}

	if opt.MIME == "" {
		indirect := typesx.DeRef(t)

		switch indirect.Kind() {
		case reflect.Slice:
			elem := indirect.Elem()
			if elem.PkgPath() == "" && elem.Kind() == reflect.Uint8 {
				opt.MIME = "plain" // bytes
			} else {
				opt.MIME = "json"
			}
		case reflect.Struct:
			// *mime/multipart.FileHeader
			if indirect.PkgPath() == "mime/multipart" && indirect.Name() == "FileHeader" {
				opt.MIME = "octet-stream"
			} else {
				opt.MIME = "json"
			}
		case reflect.Map, reflect.Array:
			opt.MIME = "json"
		default:
			opt.MIME = "plain"
		}

		if _, ok := typesx.EncodingTextMarshalerTypeReplacer(t); ok {
			opt.MIME = "plain"
		}
	}

	if ct, ok := f.set[opt.MIME]; ok {
		tsf, err := ct.New(ContextWithFactory(ctx, f), t)
		if err != nil {
			return nil, err
		}
		f.cache.Store(key, tsf)
		return tsf, nil
	}
	return nil, errors.Errorf("fmt %s is not supported for content transformer", key)
}

// Transformers returns all tsfms registered to DefaultFactory
func Transformers() (ret []string) {
	names := make(map[string]bool)
	for _, tf := range DefaultFactory.set {
		name := tf.Names()[0]
		if s, ok := tf.(CanString); ok {
			name = s.String()
		}
		if _, ok := names[name]; !ok {
			names[name] = true
			ret = append(ret, name)
		}
	}
	sort.Strings(ret)
	return ret
}
