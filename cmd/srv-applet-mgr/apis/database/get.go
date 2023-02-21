package database

import (
	"context"
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/types"
)

type GetUserDB struct {
	httpx.MethodGet
}

func (r *GetUserDB) Path() string {
	return "/user_data/"
}

func (r *GetUserDB) Output(ctx context.Context) (interface{}, error) {
	db, ok := types.WasmDBExecutorFromContext(ctx)
	if !ok {
		return nil, errors.New("fail to load db")
	}
	endpoint, err := db.ReadOnlyUser()
	if err != nil {
		return nil, err
	}
	return endpoint, nil
}
