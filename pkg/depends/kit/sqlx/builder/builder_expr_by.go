package builder

import "context"

type by struct {
	build func(context.Context) *Ex
}

func (b *by) IsNil() bool {
	return b == nil || b.build == nil
}

func (b *by) Ex(ctx context.Context) *Ex {
	return b.build(ctx)
}
