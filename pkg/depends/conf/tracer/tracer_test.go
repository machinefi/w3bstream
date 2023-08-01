package tracer_test

import (
	"context"
	"testing"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/conf/tracer"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/datatypes"
)

func do1(ctx context.Context, tr trace.Tracer) {
	_, span := tr.Start(ctx, "do1", trace.WithTimestamp(time.Now()))
	defer span.End()

	time.Sleep(100 * time.Millisecond)
}

func do2(ctx context.Context, tr trace.Tracer) {
	_, span := tr.Start(ctx, "do2", trace.WithTimestamp(time.Now()))
	defer span.End()

	time.Sleep(100 * time.Millisecond)
}

func TestConfig(t *testing.T) {
	c := tracer.Config{
		GrpcEndpoint: types.Endpoint{
			Scheme:   "http",
			Hostname: "localhost",
			Port:     4317,
		},
		HttpEndpoint: types.Endpoint{
			Scheme:   "http",
			Hostname: "localhost",
			Port:     4318,
		},
		ServiceVersion: "1.0.0",
		InstanceID:     "unique",
		DebugMode:      datatypes.TRUE,
	}
	// use grpc endpoint
	cGRPC := c
	cGRPC.HttpEndpoint = types.Endpoint{}
	cGRPC.ServiceName = "test_config_grpc"

	// use http endpoint
	cHTTP := c
	cHTTP.GrpcEndpoint = types.Endpoint{}
	cHTTP.ServiceName = "test_config_http"

	for _, c := range []*tracer.Config{&cGRPC, &cHTTP} {
		c.SetDefault()
		err := c.Init()
		if err != nil {
			t.Log(err)
			return
		}

		tr := otel.Tracer(c.ServiceName)
		ctx, span := tr.Start(context.Background(), "TestConfig", trace.WithTimestamp(time.Now()))

		do1(ctx, tr)
		do2(ctx, tr)

		span.End()

		_ = c.Shutdown(context.Background())
	}
}
