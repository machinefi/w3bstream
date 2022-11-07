package metax

import "context"

type Ctx struct {
	ctx context.Context
}

func Context(ctx context.Context) Ctx { return Ctx{ctx: ctx} }

func (c Ctx) WithContext(ctx context.Context) Ctx { c.ctx = ctx; return c }

func (c Ctx) Context() context.Context {
	if c.ctx == nil {
		return context.Background()
	}
	return c.ctx
}
