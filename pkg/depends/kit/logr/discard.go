package logr

import "context"

func Discard() Logger { return &discard{} }

type discard struct{}

func (d *discard) WithValues(...interface{}) Logger { return d }

func (d *discard) Start(ctx context.Context, _ string, _ ...interface{}) (context.Context, Logger) {
	return ctx, d
}

func (discard) End() {}

func (discard) Debug(string, ...interface{}) {}

func (discard) Info(string, ...interface{}) {}

func (discard) Warn(error) {}

func (discard) Error(error) {}
