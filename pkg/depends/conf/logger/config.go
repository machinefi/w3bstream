package logger

import (
	"os"

	"golang.org/x/exp/slog"

	"github.com/machinefi/w3bstream/pkg/depends/base/consts"
	"github.com/machinefi/w3bstream/pkg/depends/kit/logr"
)

type Config struct {
	Service string
	Version string
	Level   logr.Level `env:""`
	Output  OutputType `env:""`
	Format  FormatType `env:""`

	init bool
}

func (c *Config) SetDefault() {
	if c.Output == 0 {
		c.Output = OUTPUT_TYPE__ALWAYS
	}
	if c.Format == 0 {
		c.Format = FORMAT_TYPE__JSON
	}
}

func (c *Config) Init() error {
	if !c.init {
		if c.Service == "" {
			c.Service = os.Getenv(consts.EnvProjectName)
		}
		if c.Version == "" {
			c.Version = os.Getenv(consts.EnvProjectVersion)
		}

		switch c.Level {
		case logr.ErrorLevel:
			gLevel = slog.LevelError
		case logr.WarnLevel:
			gLevel = slog.LevelWarn
		case logr.InfoLevel:
			gLevel = slog.LevelInfo
		case logr.DebugLevel:
			gLevel = slog.LevelDebug
		}

		gOutput = c.Output

		var handler slog.Handler
		/*
			skip :[
				1. runtime.Callers,
				2. pkg/depends/conf/logger.(*customJsonHandler).Handle
				3. golang.org/x/exp/slog.(*Logger).logAttrs
				4. golang.org/x/exp/slog.(*Logger).LogAttrs
				5. pkg/depends/conf/logger.(*stdoutSpanExporter).ExportSpans
				6. pkg/depends/conf/logger.(*spanMapExporter).ExportSpans
				7. go.opentelemetry.io/otel/sdk/trace.(*simpleSpanProcessor).OnEnd
				8. go.opentelemetry.io/otel/sdk/trace.(*recordingSpan).End
				9. pkg/depends/conf/logger.(*spanLogger).End
			]
		*/

		switch c.Format {
		case FORMAT_TYPE__TEXT:
			handler = NewTextHandler(0)
		default:
			handler = NewJsonHandler(0)
		}
		gStdLogger = slog.New(handler)

		c.init = true
	}
	return nil
}

var (
	// gOutput global output to trace provider option, depends on span's status
	gOutput OutputType
	// gLevel global log level
	gLevel slog.Level
	// gStdLogger just used for stdoutSpanExporter
	gStdLogger *slog.Logger
)
