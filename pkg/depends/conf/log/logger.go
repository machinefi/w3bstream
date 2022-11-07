package log

import (
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
)

type Logger interface {
	// Start to start span for tracing
	//
	// 	ctx log = log.Start(ctx, "SpanName")
	// 	defer log.End()
	//
	Start(context.Context, string, ...interface{}) (context.Context, Logger)
	// End to end span
	End()

	// WithValues key value pairs
	WithValues(keyAndValues ...interface{}) Logger

	Trace(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(err error)
	Error(err error)
	Fatal(err error)
	Panic(err error)
}

type keyLogger struct{}

func WithLogger(ctx context.Context, l Logger) context.Context {
	return contextx.WithValue(ctx, keyLogger{}, l)
}

func FromContext(ctx context.Context) Logger {
	if v, ok := ctx.Value(keyLogger{}).(Logger); ok {
		return v
	}
	return Discard()
}

func Start(ctx context.Context, name string, kvs ...interface{}) (context.Context, Logger) {
	return FromContext(ctx).Start(ctx, name, kvs...)
}

type LevelSetter interface {
	SetLevel(lvl Level) Logger
}

// Level type
type Level uint32

const (
	// PanicLevel level, the highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel Level = iota
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel
)

// ParseLevel takes a string level and returns the Logrus log level constant.
func ParseLevel(lvl string) (Level, error) {
	switch strings.ToLower(lvl) {
	case "panic":
		return PanicLevel, nil
	case "fatal":
		return FatalLevel, nil
	case "error":
		return ErrorLevel, nil
	case "warn", "warning":
		return WarnLevel, nil
	case "info":
		return InfoLevel, nil
	case "debug":
		return DebugLevel, nil
	case "trace":
		return TraceLevel, nil
	}

	var l Level
	return l, fmt.Errorf("not a valid level: %q", lvl)
}

func (lvl Level) String() string {
	if b, err := lvl.MarshalText(); err == nil {
		return string(b)
	} else {
		return "unknown"
	}
}

func (lvl *Level) UnmarshalText(text []byte) error {
	l, err := ParseLevel(string(text))
	if err != nil {
		return err
	}
	*lvl = l
	return nil
}

func (lvl Level) MarshalText() ([]byte, error) {
	switch lvl {
	case TraceLevel:
		return []byte("trace"), nil
	case DebugLevel:
		return []byte("debug"), nil
	case InfoLevel:
		return []byte("info"), nil
	case WarnLevel:
		return []byte("warning"), nil
	case ErrorLevel:
		return []byte("error"), nil
	case FatalLevel:
		return []byte("fatal"), nil
	case PanicLevel:
		return []byte("panic"), nil
	}

	return nil, fmt.Errorf("not a valid level %d", lvl)
}

func (lvl Level) LogrusLogLevel() logrus.Level { return logrus.Level(lvl) }
