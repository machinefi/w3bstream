package contextx

import (
	"context"
	"reflect"
)

// WithValue like context.WithValue but faster
func WithValue(parent context.Context, k, v interface{}) context.Context {
	if parent == nil {
		panic("parent is nil")
	}
	if k == nil {
		panic("key is nil")
	}
	return &kvc{parent, k, v}
}

type kvc struct {
	context.Context
	k, v interface{}
}

func (c *kvc) String() string {
	return nameof(c.Context) +
		".WithValue(type " + reflect.TypeOf(c.k).String() +
		", val" + stringify(c.v) + ")"
}

func (c *kvc) Value(k interface{}) interface{} {
	if c.k == k {
		return c.v
	}
	return c.Context.Value(k)
}

type stringer interface{ String() string }

func nameof(c context.Context) string {
	if str, ok := c.(stringer); ok {
		return str.String()
	}
	return reflect.TypeOf(c).String()
}

func stringify(v interface{}) string {
	switch s := v.(type) {
	case stringer:
		return s.String()
	case string:
		return s
	}
	return "<not Stringer>"
}
