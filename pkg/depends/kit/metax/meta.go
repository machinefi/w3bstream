package metax

import (
	"context"
	"net/url"
	"strings"

	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
)

type key struct{}

func ContextWith(ctx context.Context, key string, values ...string) context.Context {
	return ContextWithMeta(ctx, Meta{key: values})
}

func ContextWithMeta(ctx context.Context, meta Meta) context.Context {
	return contextx.WithValue(ctx, key{}, GetMetaFrom(ctx).Merge(meta))
}

func GetMetaFrom(ctx context.Context) Meta {
	if m, ok := ctx.Value(key{}).(Meta); ok {
		return m
	}
	return Meta{}
}

func ParseMeta(query string) Meta {
	if !strings.Contains(query, "=") {
		return Meta{
			"_id": []string{query},
		}
	}
	values, err := url.ParseQuery(query)
	if err == nil {
		return Meta(values)
	}
	return Meta{}
}

type Meta map[string][]string

func (m Meta) Merge(metas ...Meta) Meta {
	meta := m.Clone()

	for _, me := range metas {
		for k, values := range me {
			if k == "" {
				continue
			}
			if k[0] == '_' {
				meta[k] = values
				continue
			}
			meta.Add(k, values...)
		}
	}

	return meta
}

func (m Meta) Clone() Meta {
	meta := Meta{}
	for k, v := range m {
		meta[k] = v
	}
	return meta
}

func (m Meta) Add(key string, values ...string) {
	m[key] = append(m[key], values...)
}

func (m Meta) Get(key string) string {
	if m == nil {
		return ""
	}
	vs := m[key]
	if len(vs) == 0 {
		return ""
	}
	return vs[0]
}

func (m Meta) With(key string, values ...string) Meta {
	meta := m.Clone()
	meta[key] = values
	return meta
}

func (m Meta) String() string {
	return url.Values(m).Encode()
}
