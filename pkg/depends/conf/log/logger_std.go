package log

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/machinefi/w3bstream/pkg/depends/x/ptrx"
)

func Std() Logger { return &std{lvl: DebugLevel} }

func StdContext(ctx context.Context) (context.Context, Logger) {
	l := Std()
	return WithLogger(ctx, l), l
}

type std struct {
	lvl   Level
	spans []string
	kvs   []interface{}
}

func (l *std) SetLevel(lvl Level) Logger {
	logger := ptrx.Ptr(*l)
	logger.lvl = lvl
	return logger
}

func (l *std) WithValues(kvs ...interface{}) Logger {
	return &std{
		lvl:   l.lvl,
		spans: l.spans,
		kvs:   append(l.kvs, kvs...),
	}
}

func (l *std) Start(ctx context.Context, name string, kvs ...interface{}) (context.Context, Logger) {
	return ctx, &std{
		lvl:   l.lvl,
		spans: append(l.spans, name),
		kvs:   append(l.kvs, kvs...),
	}
}

func (l *std) End() {
	if len(l.spans) != 0 {
		l.spans = l.spans[0 : len(l.spans)-1]
	}
}

func (l *std) Trace(format string, args ...interface{}) {
	if TraceLevel <= l.lvl {
		logrus.WithFields(KeyValues(l.kvs...)).Tracef(format, args...)
	}
}

func (l *std) Debug(format string, args ...interface{}) {
	if DebugLevel <= l.lvl {
		logrus.WithFields(KeyValues(l.kvs...)).Debugf(format, args...)
	}
}

func (l *std) Info(format string, args ...interface{}) {
	if InfoLevel <= l.lvl {
		logrus.WithFields(KeyValues(l.kvs...)).Infof(format, args...)
	}
}

func (l *std) Warn(err error) {
	if WarnLevel <= l.lvl {
		logrus.WithFields(KeyValues(l.kvs...)).Warn(err)
	}
}

func (l *std) Error(err error) {
	if ErrorLevel <= l.lvl {
		logrus.WithFields(KeyValues(l.kvs...)).Error(err)
	}
}

func (l *std) Fatal(err error) {
	if FatalLevel <= l.lvl {
		logrus.WithFields(KeyValues(l.kvs...)).Fatal(err)
	}
}

func (l *std) Panic(err error) {
	if PanicLevel < l.lvl {
		logrus.WithFields(KeyValues(l.kvs...)).Panic(err)
	}
}

func KeyValues(kvs ...interface{}) logrus.Fields {
	if len(kvs)%2 != 0 {
		return nil
	}

	fields := logrus.Fields{}

	for i := 0; i < len(kvs); i += 2 {
		fields[fmt.Sprintf("%v", kvs[i])] = kvs[i+1]
	}

	return fields
}
