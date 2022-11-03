package builder

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
)

type Toggles map[string]bool

const (
	ToggleNeedAutoAlias = "NeedAlias"
	ToggleMultiTable    = "MultiTable"
	ToggleUseValues     = "UseValues"
)

func (toggles Toggles) Merge(follows Toggles) Toggles {
	final := Toggles{}
	for k, v := range toggles {
		if v {
			final[k] = true
		}
	}
	for k, v := range follows {
		if v {
			final[k] = true
		} else {
			delete(final, k)
		}
	}
	return final
}

func (toggles Toggles) Is(k string) bool {
	v, ok := toggles[k]
	return ok && v
}

type ctxKeyToggles struct{}

func ContextWithToggles(ctx context.Context, toggles Toggles) context.Context {
	return contextx.WithValue(
		ctx,
		ctxKeyToggles{},
		TogglesFromContext(ctx).Merge(toggles),
	)
}

func ContextWithToggleMultiTable(ctx context.Context, v bool) context.Context {
	return ContextWithToggles(ctx, Toggles{ToggleMultiTable: v})
}

func ContextWithToggleUseValues(ctx context.Context, v bool) context.Context {
	return ContextWithToggles(ctx, Toggles{ToggleUseValues: v})
}

func ContextWithToggleNeedAutoAlias(ctx context.Context, v bool) context.Context {
	return ContextWithToggles(ctx, Toggles{ToggleNeedAutoAlias: v})
}

func TogglesFromContext(ctx context.Context) Toggles {
	if ctx == nil {
		return Toggles{}
	}
	if toggles, ok := ctx.Value(ctxKeyToggles{}).(Toggles); ok {
		return toggles
	}
	return Toggles{}
}
