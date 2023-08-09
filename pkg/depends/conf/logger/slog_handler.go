package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"golang.org/x/exp/slog"
)

func NewJsonHandler(skip int) *customJsonHandler {
	return &customJsonHandler{
		skip: skip,
		JSONHandler: slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			AddSource:   skip > 0,
			Level:       gLevel,
			ReplaceAttr: CustomReplacer,
		}),
	}
}

type customJsonHandler struct {
	skip int
	*slog.JSONHandler
}

func (h *customJsonHandler) Handle(ctx context.Context, r slog.Record) error {
	var pcs [1]uintptr
	runtime.Callers(h.skip, pcs[:])
	r.PC = pcs[0]
	return h.JSONHandler.Handle(ctx, r)
}

func (h *customJsonHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &customJsonHandler{
		skip:        h.skip,
		JSONHandler: h.JSONHandler.WithAttrs(attrs).(*slog.JSONHandler),
	}
}

func (h *customJsonHandler) WithGroup(group string) slog.Handler {
	return &customJsonHandler{
		skip:        h.skip,
		JSONHandler: h.JSONHandler.WithGroup(group).(*slog.JSONHandler),
	}
}

func (h *customJsonHandler) Enabled(_ context.Context, lv slog.Level) bool {
	return lv >= gLevel
}

func NewTextHandler(skip int) *customTextHandler {
	return &customTextHandler{
		skip: skip,
		TextHandler: slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			AddSource:   skip > 0,
			Level:       gLevel,
			ReplaceAttr: CustomReplacer,
		}),
	}
}

type customTextHandler struct {
	skip int
	*slog.TextHandler
}

func (h *customTextHandler) Handle(ctx context.Context, r slog.Record) error {
	var pcs [1]uintptr
	runtime.Callers(h.skip, pcs[:])
	r.PC = pcs[0]
	return h.TextHandler.Handle(ctx, r)
}

func (h *customTextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &customTextHandler{
		skip:        h.skip,
		TextHandler: h.TextHandler.WithAttrs(attrs).(*slog.TextHandler),
	}
}

func (h *customTextHandler) WithGroup(group string) slog.Handler {
	return &customTextHandler{
		skip:        h.skip,
		TextHandler: h.TextHandler.WithGroup(group).(*slog.TextHandler),
	}
}

func (h *customTextHandler) Enabled(_ context.Context, lv slog.Level) bool {
	return lv >= gLevel
}

func CustomReplacer(groups []string, a slog.Attr) slog.Attr {
	// time format
	v := a.Value
	if a.Value.Kind() == slog.KindTime {
		a.Value = slog.StringValue(v.Time().Format("20060102-150405.000Z0700"))
	}

	// replace time key format
	if a.Key == slog.TimeKey {
		a.Key = "@ts"
		return a
	}

	// replace level key
	if a.Key == slog.LevelKey {
		val := ""
		switch v := a.Value.Any().(slog.Level); v {
		case slog.LevelDebug:
			val = "deb"
		case slog.LevelInfo:
			val = "inf"
		case slog.LevelWarn:
			val = "wrn"
		default:
			val = "err"
		}
		return slog.String("@lv", val)
	}

	// replace stack info
	if a.Key == slog.SourceKey {
		s := a.Value.Any().(*slog.Source)
		s.File = filepath.Base(s.File)
		return slog.String("@src", fmt.Sprintf("%s:%d", s.File, s.Line))
	}

	return a
}
