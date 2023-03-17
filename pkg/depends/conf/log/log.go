package log

import (
	"context"
	"os"
	"runtime"
	"strconv"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/sdk/trace"

	"github.com/machinefi/w3bstream/pkg/depends/base/consts"
	"github.com/machinefi/w3bstream/pkg/depends/kit/metax"
)

type Log struct {
	Name         string
	Level        Level            `env:""`
	Output       LoggerOutputType `env:""`
	Format       LoggerFormatType
	Exporter     trace.SpanExporter `env:"-"`
	ReportCaller bool
}

func (l *Log) SetDefault() {
	if l.Level == 0 {
		l.Level = DebugLevel
	}
	if l.Output == 0 {
		l.Output = LOGGER_OUTPUT_TYPE__ALWAYS
	}
	if l.Format == 0 {
		l.Format = LOGGER_FORMAT_TYPE__JSON
	}
	if l.Name == "" {
		l.Name = "unknown"
		if v := os.Getenv(consts.EnvProjectName); v != "" {
			l.Name = v
		}
	}
}

func (l *Log) InitLogrus() {
	if l.Format == LOGGER_FORMAT_TYPE__JSON {
		logrus.SetFormatter(JsonFormatter)
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	}

	logrus.SetLevel(l.Level.LogrusLogLevel())
	logrus.SetReportCaller(l.ReportCaller)
	// TODO add hook with goid meta logrus.AddHook(goid.Default)
	logrus.AddHook(&ProjectAndMetaHook{l.Name})

	logrus.SetOutput(os.Stdout)
}

func (l *Log) InitSpanLog() {
	if l.Exporter == nil {
		return
	}
	if err := InstallPipeline(l.Output, l.Format, l.Exporter); err != nil {
		panic(err)
	}
}

func (l *Log) Init() {
	l.InitLogrus()
	l.InitSpanLog()
}

type ProjectAndMetaHook struct {
	Name string
}

func (h *ProjectAndMetaHook) Fire(entry *logrus.Entry) error {
	ctx := entry.Context
	if ctx == nil {
		ctx = context.Background()
	}
	meta := metax.GetMetaFrom(ctx)
	if entry.Data["@prj"] == nil {
		entry.Data["@prj"] = h.Name
	}
	for k, v := range meta {
		entry.Data["meta."+k] = v
	}
	return nil
}

func (h *ProjectAndMetaHook) Levels() []logrus.Level { return logrus.AllLevels }

var (
	project       = "unknown"
	JsonFormatter = &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyLevel: "@lv",
			logrus.FieldKeyTime:  "@ts",
			logrus.FieldKeyFunc:  "@fn",
			logrus.FieldKeyFile:  "@fl",
		},
		CallerPrettyfier: func(f *runtime.Frame) (fn string, file string) {
			return f.Function + " line:" + strconv.FormatInt(int64(f.Line), 10), ""
		},
		TimestampFormat: "20060102-150405.000Z07:00",
	}
)

func init() {
	if v := os.Getenv(consts.EnvProjectName); v != "" {
		project = v
		if version := os.Getenv(consts.EnvProjectVersion); version != "" {
			project = project + "@" + version
		}
	}
	logrus.SetFormatter(JsonFormatter)
	logrus.SetReportCaller(true)
}
