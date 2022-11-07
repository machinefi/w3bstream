package goid

import (
	"sync"

	"github.com/machinefi/w3bstream/pkg/depends/kit/metax"
)

var Default = &Meta{}

type Meta struct{ m sync.Map }

func (m *Meta) Clear() {
	// m.m.Delete(runtime.GoID())
}

func (m *Meta) Get() metax.Meta {
	// if logID, ok := m.m.Load(runtime.GoID()); ok {
	// 	return logID.(metax.Meta)
	// }
	return metax.Meta{}
}

func (m *Meta) Set(meta metax.Meta) {
	// m.m.Store(runtime.GoID(), meta)
}

func (m *Meta) With(cb func(), metas ...metax.Meta) func() {
	meta := metax.Meta{}

	if len(metas) == 0 {
		meta = m.Get()
	} else {
		meta = meta.Merge(metas...)
	}

	return func() {
		m.Set(meta)
		defer m.Clear()
		cb()
	}
}

func (m *Meta) All() map[int64]metax.Meta {
	results := map[int64]metax.Meta{}

	m.m.Range(func(key, value interface{}) bool {
		results[key.(int64)] = value.(metax.Meta)
		return true
	})

	return results
}
