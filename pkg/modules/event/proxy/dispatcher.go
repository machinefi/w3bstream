// event dispatch

package proxy

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/iotexproject/w3bstream/pkg/modules/deploy"
	"github.com/iotexproject/w3bstream/pkg/modules/vm"
	"github.com/iotexproject/w3bstream/pkg/types"
	"github.com/iotexproject/w3bstream/pkg/types/wasm"
)

func dispatch(ctx context.Context, e Event) ([]byte, error) {
	l := types.MustLoggerFromContext(ctx)

	_, l = l.Start(ctx, "EventDispatch")

	ins, err := deploy.GetInstanceByAppletID(ctx, e.Meta().AppletID)
	if err != nil {
		l.Error(err)
		return nil, err
	}
	if len(ins) == 0 {
		err = errors.Errorf("applet not found")
		l.Warn(err)
		return nil, err
	}
	consumer := vm.GetConsumer(ins[0].InstanceID)
	if consumer == nil {
		err = errors.Errorf("instance not found")
		l.Error(err)
		return nil, err
	}
	cctx, _ := context.WithTimeout(ctx, time.Second*3)
	res, code, err := consumer.HandleEvent(cctx, e.Meta().Handler, e.Raw())
	if err != nil {
		l.Error(err)
		return nil, err
	}
	if code != wasm.ResultStatusCode_OK {
		err = errors.Errorf("wasm failed, error code %v", code)
		l.Error(err)
		return nil, err
	}
	return res, nil
}
