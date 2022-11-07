package http

import (
	"context"
	"reflect"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/depends/x/reflectx"
)

type LivenessChecker interface {
	LivenessCheck() map[string]string
}

var (
	RTypeLivenessChecker = reflect.TypeOf((*LivenessChecker)(nil)).Elem()
)

func RegisterCheckerBy(vs ...interface{}) {
	for _, v := range vs {
		rv := reflectx.Indirect(reflect.ValueOf(v))
		rt := rv.Type()

		if chk, ok := v.(LivenessChecker); ok {
			RegisterChecker(rt.Name(), chk)
			continue
		}

		if rt.Kind() != reflect.Struct {
			continue
		}

		for i := 0; i < rv.NumField(); i++ {
			fv := rv.Field(i)
			ft := rt.Field(i)

			if !ft.IsExported() {
				continue
			}
			if chk, ok := fv.Interface().(LivenessChecker); ok {
				RegisterChecker(ft.Name, chk)
			}
		}
	}
}

func RegisterChecker(k string, chk LivenessChecker) { checkers[k] = chk }

func ResetRegistered() { checkers = LivenessCheckers{} }

type LivenessCheckers map[string]LivenessChecker

var checkers = LivenessCheckers{}

func (cs LivenessCheckers) Statuses() map[string]string {
	m := map[string]string{}

	for name := range cs {
		if cs[name] != nil {
			entry := cs[name].LivenessCheck()
			for key, v := range entry {
				if key != "" {
					m[name+"/"+key] = v
				} else {
					m[name] = v
				}
			}
		}
	}

	return m
}

var LivenessRouter = kit.NewRouter(&Liveness{})

type Liveness struct{ httpx.MethodGet }

func (Liveness) Path() string { return "/liveness" }

func (Liveness) Output(ctx context.Context) (interface{}, error) {
	return checkers.Statuses(), nil
}

// type LivenessStatuses map[string]string
