package id

import (
	"context"
	"net"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/base/types/snowflake_id"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/must"
)

var (
	start, _ = time.Parse(time.RFC3339, "2022-10-13T22:10:24+08:00")
	sff      = snowflake_id.NewSnowflakeFactory(12, 10, 1, start.Local())
)

type Generator interface {
	ID() (uint64, error)
}

func FromIP(ip net.IP) (Generator, error) {
	return sff.NewSnowflake(snowflake_id.WorkerIDFromIP(ip))
}

func FromLocalIP() (Generator, error) {
	return sff.NewSnowflake(snowflake_id.WorkerIDFromLocalIP())
}

type keyGenerator struct{}

func WithGenerator(ctx context.Context, g Generator) context.Context {
	return contextx.WithValue(ctx, keyGenerator{}, g)
}

func WithGeneratorContext(g Generator) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, keyGenerator{}, g)
	}
}

func GeneratorFromContext(ctx context.Context) (Generator, bool) {
	g, ok := ctx.Value(keyGenerator{}).(Generator)
	return g, ok
}

func MustGeneratorFromContext(ctx context.Context) Generator {
	g, ok := ctx.Value(keyGenerator{}).(Generator)
	must.BeTrue(ok)
	return g
}

type SFIDGenerator interface {
	MustGenSFID() types.SFID
	MustGenSFIDs(n int) types.SFIDs
}

type keySFIDGenerator struct{}

func WithSFIDGenerator(ctx context.Context, g SFIDGenerator) context.Context {
	return contextx.WithValue(ctx, keySFIDGenerator{}, g)
}

func WithSFIDGeneratorContext(g SFIDGenerator) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, keySFIDGenerator{}, g)
	}
}

func SFIDGeneratorFromContext(ctx context.Context) (SFIDGenerator, bool) {
	g, ok := ctx.Value(keySFIDGenerator{}).(SFIDGenerator)
	return g, ok
}

func MustSFIDGeneratorFromContext(ctx context.Context) SFIDGenerator {
	g, ok := ctx.Value(keySFIDGenerator{}).(SFIDGenerator)
	must.BeTrue(ok)
	return g
}
