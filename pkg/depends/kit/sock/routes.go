package sock

import (
	"sync"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sche"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sock/sock_msg"
)

type Handler func(*Event)

func HandlerFunc(h Handler, ev *Event) sche.Fn { return func() { h(ev) } }

type Job func()

type Routes struct {
	mu *sync.Mutex
	v  map[sock_msg.Type][]Handler
}

func NewRoutes() *Routes {
	return &Routes{
		mu: &sync.Mutex{},
		v:  make(map[sock_msg.Type][]Handler),
	}
}

func (r *Routes) Register(t sock_msg.Type, fns ...Handler) {
	if r == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.v[t] = append(r.v[t], fns...)
}

func (r *Routes) Handlers(t sock_msg.Type) []Handler {
	if r == nil {
		return nil
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.v[t]
}
