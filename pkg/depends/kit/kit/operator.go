package kit

import (
	"fmt"
	"net/url"
	"reflect"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/x/reflectx"
)

func NewOperatorFactory(op Operator, last bool) *OperatorFactory {
	opType := reflectx.DeRef(reflect.TypeOf(op))
	if opType.Kind() != reflect.Struct {
		panic(fmt.Errorf("operator must be a struct type, got %#v", op))
	}

	meta := &OperatorFactory{}
	meta.IsLast = last

	meta.Operator = op

	if _, ok := op.(OperatorWithoutOutput); ok {
		meta.NoOutput = true
	}

	meta.Type = reflectx.DeRef(reflect.TypeOf(op))

	if operatorWithParams, ok := op.(OperatorWithParams); ok {
		meta.Params = operatorWithParams.OperatorParams()
	}

	if !meta.IsLast {
		if ctxKey, ok := op.(ContextProvider); ok {
			meta.ContextKey = ctxKey.ContextKey()
		} else {
			if ctxKey, ok := op.(oldContextProvider); ok {
				meta.ContextKey = ctxKey.ContextKey()
			} else {
				meta.ContextKey = meta.Type.String()
			}
		}
	}

	return meta
}

type oldContextProvider interface {
	ContextKey() string
}

// TODO remove ContextKey?

type OperatorFactory struct {
	Type       reflect.Type
	ContextKey interface{}
	NoOutput   bool
	Params     url.Values
	IsLast     bool
	Operator   Operator
}

func (o *OperatorFactory) String() string {
	if o.Params != nil {
		return o.Type.String() + "?" + o.Params.Encode()
	}
	return o.Type.String()
}

func (o *OperatorFactory) New() Operator {
	var op Operator

	if newer, ok := o.Operator.(OperatorNewer); ok {
		op = newer.New()
	} else {
		op = reflect.New(o.Type).Interface().(Operator)
	}

	if init, ok := op.(OperatorInit); ok {
		init.InitFrom(o.Operator)
	}

	if setter, ok := op.(types.DefaultSetter); ok {
		setter.SetDefault()
	}

	return op
}

type EmptyOperator struct {
	OperatorWithoutOutput
}
