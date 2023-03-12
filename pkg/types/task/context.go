package task

//
//import (
//	"context"
//	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
//	"github.com/machinefi/w3bstream/pkg/depends/x/misc/must"
//	"github.com/machinefi/w3bstream/pkg/modules/job"
//)
//
//type (
//	CtxDispatcher struct{}
//)
//
//func WithDispatcher(ctx context.Context, v job.Dispatcher) context.Context {
//	return contextx.WithValue(ctx, CtxDispatcher{}, v)
//}
//
//func WithDispatcherContext(v job.Dispatcher) contextx.WithContext {
//	return func(ctx context.Context) context.Context {
//		return contextx.WithValue(ctx, CtxDispatcher{}, v)
//	}
//}
//
//func DispatcherFromContext(ctx context.Context) (job.Dispatcher, bool) {
//	v, ok := ctx.Value(CtxDispatcher{}).(job.Dispatcher)
//	return v, ok
//}
//
//func MustDispatcherFromContext(ctx context.Context) job.Dispatcher {
//	v, ok := DispatcherFromContext(ctx)
//	must.BeTrue(ok)
//	return v
//}
