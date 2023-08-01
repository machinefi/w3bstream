package tracer

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.16.0"
	"google.golang.org/grpc/credentials"

	"github.com/machinefi/w3bstream/pkg/depends/base/consts"
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/conf/logger"
	conftls "github.com/machinefi/w3bstream/pkg/depends/conf/tls"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
)

type Config struct {
	// GrpcEndpoint provider grpc endpoint
	GrpcEndpoint types.Endpoint `env:""`
	// HttpEndpoint provider http endpoint
	HttpEndpoint types.Endpoint `env:""`
	// ServiceName service name, default use env `PRJ_NAME`
	ServiceName string `env:""`
	// ServiceVersion service version, default use env `PRJ_VERSION`
	ServiceVersion string `env:""`
	// InstanceID the unique id of current service, default gen uuid when trace config init
	InstanceID string `env:""`
	// DebugMode if enable debug mode
	DebugMode datatypes.Bool `env:""`
	// TLS config for connecting provider Endpoint
	TLS *conftls.X509KeyPair

	provider *trace.TracerProvider
}

func (c *Config) IsZero() bool {
	return c.GrpcEndpoint.IsZero() && c.HttpEndpoint.IsZero()
}

func (c *Config) SetDefault() {
	if !c.GrpcEndpoint.IsZero() && c.GrpcEndpoint.Port == 0 {
		c.GrpcEndpoint.Port = 4317
	}
	if !c.HttpEndpoint.IsZero() && c.HttpEndpoint.Port == 0 {
		c.HttpEndpoint.Port = 4318
	}
	if c.DebugMode == 0 {
		c.DebugMode = datatypes.FALSE
	}
}

func (c *Config) grpcExporter() (*otlptrace.Exporter, error) {
	options := []otlptracegrpc.Option{otlptracegrpc.WithEndpoint(c.GrpcEndpoint.Host())}
	if c.GrpcEndpoint.IsTLS() {
		if !c.TLS.IsZero() {
			options = append(options, otlptracegrpc.WithTLSCredentials(
				credentials.NewClientTLSFromCert(c.TLS.TLSConfig().RootCAs, ""),
			))
		}
	} else {
		options = append(options, otlptracegrpc.WithInsecure())
	}

	return otlptracegrpc.New(context.Background(), options...)
}

func (c *Config) httpExporter() (*otlptrace.Exporter, error) {
	options := []otlptracehttp.Option{otlptracehttp.WithEndpoint(c.HttpEndpoint.Host())}
	if c.GrpcEndpoint.IsTLS() {
		if !c.TLS.IsZero() {
			options = append(options, otlptracehttp.WithTLSClientConfig(c.TLS.TLSConfig()))
		}
	} else {
		options = append(options, otlptracehttp.WithInsecure())
	}

	return otlptracehttp.New(context.Background(), options...)
}

func (c *Config) Init() error {
	if c.IsZero() {
		return nil // return nil and using global.defaultTracerValue
	}

	if c.ServiceName == "" {
		c.ServiceName = os.Getenv(consts.EnvProjectName)
	}
	if c.ServiceVersion == "" {
		c.ServiceVersion = os.Getenv(consts.EnvProjectVersion)
	}
	if c.InstanceID == "" {
		c.InstanceID = uuid.NewString()
	}

	var (
		exp trace.SpanExporter
		err error
		ep  types.Endpoint
	)

	if !c.GrpcEndpoint.IsZero() {
		ep = c.GrpcEndpoint
		exp, err = c.grpcExporter()
	} else {
		ep = c.HttpEndpoint
		exp, err = c.httpExporter()
	}
	if err != nil {
		return errors.Errorf("new exporter failed: ep[%v] err[%v]", ep, err)
	}

	// resource
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(c.ServiceName),
		semconv.ServiceVersionKey.String(c.ServiceVersion),
		attribute.String("instance", c.InstanceID),
	)
	if err != nil {
		return errors.Errorf("new otlp resource failed: %v", err)
	}

	options := []trace.TracerProviderOption{
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithResource(res),
		trace.WithSyncer(logger.WithSpanMapExporter(logger.OutputFilter())(logger.StdoutSpanExporter())),
	}
	if c.DebugMode == datatypes.TRUE {
		options = append(options,
			trace.WithBatcher(logger.WithSpanMapExporter(logger.OutputFilter())(exp)),
		)
	} else {
		options = append(options,
			trace.WithBatcher(logger.WithSpanMapExporter(logger.OutputFilter())(logger.WithErrIgnoreExporter()(exp))),
		)
	}

	c.provider = trace.NewTracerProvider(options...)

	// set global trace provider
	otel.SetTracerProvider(c.provider)
	otel.SetTextMapPropagator(propagation.Baggage{})

	log.Printf("Trace provider for service `%s@%s` initialized\n", c.ServiceName, c.ServiceVersion)
	go func() {
		stopCh := make(chan os.Signal, 1)
		signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)
		<-stopCh
		err := c.Shutdown(context.Background())
		log.Printf("Trace provider for service `%s@%s` shutdown: %v\n", c.ServiceName, c.ServiceVersion, err)
	}()

	return nil
}

func (c *Config) Shutdown(ctx context.Context) error {
	err := c.provider.Shutdown(ctx)
	log.Printf("Trace provider for service `%s@%s` shutdown: %v\n", c.ServiceName, c.ServiceVersion, err)
	return err
}
