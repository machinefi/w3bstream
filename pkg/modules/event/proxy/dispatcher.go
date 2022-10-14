// event dispatch

package proxy

import (
	"context"
	"fmt"

	"github.com/iotexproject/w3bstream/pkg/modules/deploy"
	"github.com/iotexproject/w3bstream/pkg/modules/event"
	"github.com/iotexproject/w3bstream/pkg/modules/vm"
	"github.com/iotexproject/w3bstream/pkg/types/wasm"
)

func dispatch(ctx context.Context, e event.Event) ([]byte, error) {
	ins, err := deploy.GetInstanceByAppletID(ctx, e.Meta().AppletID)
	if err != nil {
		return nil, err
	}
	if len(ins) == 0 {
		return nil, fmt.Errorf("applet not found")
	}
	consumer := vm.GetConsumer(ins[0].InstanceID.String())
	if consumer == nil {
		return nil, fmt.Errorf("instance not found")
	}
	res, code := consumer.HandleEvent(e.Meta().Handler, e.Raw())
	if code != wasm.ResultStatusCode_OK {
		return nil, fmt.Errorf("wasm failed, error code %v", code)
	}
	return res, nil
}
