package slog

import (
	"context"

	"golang.org/x/exp/slog"
)

func Default() *slog.Logger {
	return slog.New(&handler{h: slog.Default().Handler()})
}

type handler struct {
	h slog.Handler
}

func (h *handler) Handle(ctx context.Context, r slog.Record) error {
	return h.h.Handle(ctx, r)
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &handler{h: h.h.WithAttrs(attrs)}
}

func (h *handler) WithGroup(name string) slog.Handler {
	return &handler{h: h.h.WithGroup(name)}
}

func (h *handler) Enabled(_ context.Context, l slog.Level) bool {
	return l >= slog.LevelDebug
}
