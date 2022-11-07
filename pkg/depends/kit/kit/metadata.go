package kit

import "net/url"

func FromMetas(metas ...Metadata) Metadata {
	m := Metadata{}
	for _, meta := range metas {
		m.Merge(meta)
	}
	return m
}

type Metadata map[string][]string

func (m Metadata) String() string {
	return url.Values(m).Encode()
}

func (m Metadata) Del(key string) {
	delete(m, key)
}

func (m Metadata) Merge(metadata Metadata) {
	for key, values := range metadata {
		m.Set(key, values...)
	}
}

func (m Metadata) Add(key, value string) {
	if values, ok := m[key]; ok {
		m[key] = append(values, value)
	} else {
		m.Set(key, value)
	}
}

func (m Metadata) Set(key string, values ...string) {
	m[key] = values
}

func (m Metadata) Has(key string) bool {
	_, ok := m[key]
	return ok
}

func (m Metadata) Get(key string) string {
	if v := m[key]; len(v) > 0 {
		return v[0]
	}
	return ""
}
