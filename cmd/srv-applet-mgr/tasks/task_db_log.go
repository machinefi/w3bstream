package tasks

import (
	"context"
	"github.com/pkg/errors"
	"reflect"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/types"
)

type DbLogger interface {
	Create(sqlx.DBExecutor) error
}

type DbLogStoring struct {
	l DbLogger
}

func (t *DbLogStoring) SetArg(v interface{}) error {
	if ctx, ok := v.(DbLogger); ok {
		t.l = ctx
		return nil
	}
	return errors.Errorf("invalid arg: %s", reflect.TypeOf(v))
}

func (t *DbLogStoring) Output(ctx context.Context) (interface{}, error) {
	return nil, t.l.Create(types.MustMgrDBExecutorFromContext(ctx))
}
