package mapx

import "sync"

type Map[K comparable, V any] struct {
	val map[K]V
	mtx *sync.RWMutex
}

func New[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{val: make(map[K]V), mtx: &sync.RWMutex{}}
}

func (m *Map[K, V]) Load(k K) (v V, ok bool) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	v, ok = m.val[k]
	return
}

func (m *Map[K, V]) Store(k K, v V) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.val[k] = v
}

func (m *Map[K, V]) StoreNX(k K, v V) bool {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	_, ok := m.val[k]
	if ok {
		return false
	}
	m.val[k] = v
	return true
}

func (m *Map[K, V]) LoadOrStore(k K, newv func() (V, error)) (V, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	v, ok := m.val[k]
	if ok {
		return v, nil
	}
	v, err := newv()
	if err != nil {
		return v, err
	}

	m.val[k] = v
	return v, nil
}

func (m *Map[K, V]) Remove(k K) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	delete(m.val, k)
}

func (m *Map[K, V]) Clear() {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.val = make(map[K]V)
}

func (m *Map[K, V]) LoadAndRemove(k K) (v V, ok bool) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	v, ok = m.val[k]
	if ok {
		defer delete(m.val, k)
	}
	return v, ok
}

func (m *Map[K, V]) Len() int {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	return len(m.val)
}

func (m *Map[K, V]) Range(f func(k K, v V) bool) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	for k, v := range m.val {
		if !f(k, v) {
			break
		}
	}
}
