package debug

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/vm"
)

// FetchInstances fetch instances states in memory
type FetchInstances struct {
	httpx.MethodGet
}

func (r *FetchInstances) Path() string { return "/instances" }

func (r *FetchInstances) Output(ctx context.Context) (interface{}, error) {
	return vm.FetchInstances(), nil
}

type GetInstance struct {
	httpx.MethodGet
	ID types.SFID `in:"path" name:"id"`
}

func (r *GetInstance) Path() string {
	return "/instance/:id"
}

func (r *GetInstance) Output(ctx context.Context) (interface{}, error) {
	state, _ := vm.GetInstanceState(r.ID)
	return state, nil
}
