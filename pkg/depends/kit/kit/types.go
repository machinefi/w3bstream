package kit

import "context"

type Client interface {
	Do(ctx context.Context, req interface{}, metas ...Metadata) Result
}

type Result interface {
	Into(v interface{}) (Metadata, error)
}

type Transport interface {
	Serve(router *Router) error
}

type Operator interface {
	Output(ctx context.Context) (interface{}, error)
}

type MiddleOperators []Operator

type WithMiddleOperators interface {
	MiddleOperators() MiddleOperators
}

type MetadataCarrier interface {
	Meta() Metadata
}

type OperatorWithParams interface {
	OperatorParams() map[string][]string
}

type OperatorWithoutOutput interface {
	Operator
	NoOutput()
}

type ContextProvider interface {
	Operator
	ContextKey() interface{}
}

type OperatorInit interface {
	InitFrom(o Operator)
}

type OperatorNewer interface {
	New() Operator
}
