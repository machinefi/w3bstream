// This is a generated source file. DO NOT EDIT
// Source: degradation_demo/operations.go

package degradation_demo

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/depends/kit/metax"
)

type DemoApi struct {
}

func (o *DemoApi) Path() string {
	return "/peer/version"
}

func (o *DemoApi) Method() string {
	return "GET"
}

func (o *DemoApi) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "degradationDemo.DemoApi")
	return cli.Do(ctx, o, metas...)
}

func (o *DemoApi) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (*DemoApiResp, kit.Metadata, error) {
	rsp := new(DemoApiResp)
	meta, err := cli.Do(ctx, o, metas...).Into(rsp)
	return rsp, meta, err
}

func (o *DemoApi) Invoke(cli kit.Client, metas ...kit.Metadata) (*DemoApiResp, kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}
