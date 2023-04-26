package migration

import (
	"context"
	"io"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
)

type key struct{}

func OutputFromContext(ctx context.Context) io.Writer {
	if opts, ok := ctx.Value(key{}).(io.Writer); ok {
		if opts != nil {
			return opts
		}
	}
	return nil
}

func WithInspectionOutput(ctx context.Context, w io.Writer) context.Context {
	return contextx.WithValue(ctx, key{}, w)
}

func MustMigrate(db sqlx.DBExecutor, w io.Writer) {
	if err := Migrate(db, w); err != nil {
		panic(err)
	}
}

func Migrate(db sqlx.DBExecutor, output io.Writer) error {
	ctx := contextx.WithValue(db.Context(), key{}, output)

	if err := db.(sqlx.Migrator).Migrate(ctx, db); err != nil {
		return err
	}
	if output == nil {
		if err := SyncEnum(db); err != nil {
			return err
		}
	}
	return nil
}
