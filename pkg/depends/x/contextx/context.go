package contextx

import "context"

type WithContext = func(ctx context.Context) context.Context

func WithContextCompose(withs ...WithContext) WithContext {
	return func(ctx context.Context) context.Context {
		for i := range withs {
			ctx = withs[i](ctx)
		}
		return ctx
	}
}
