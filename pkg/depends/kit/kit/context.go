package kit

import "context"

type ContextWith = func(ctx context.Context) context.Context

func ComposeContextWith(withs ...ContextWith) ContextWith {
	return func(ctx context.Context) context.Context {
		for i := range withs {
			ctx = withs[i](ctx)
		}
		return ctx
	}
}
