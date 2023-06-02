package robot_notifier

import (
	"bytes"
	"context"
	"encoding"
	"encoding/json"
	"io"
	"net/http"

	"github.com/machinefi/w3bstream/pkg/types"
)

func Push(ctx context.Context, data interface{}, rspHookFn ...func([]byte) error) (err error) {
	notifier, ok := types.RobotNotifierConfigFromContext(ctx)
	if !ok || notifier.IsZero() {
		return nil
	}

	var (
		msg  []byte
		body []byte
	)
	switch v := data.(type) {
	case []byte:
		msg = v
	case string:
		msg = []byte(v)
	case encoding.TextMarshaler:
		msg, err = v.MarshalText()
	default:
		msg, err = json.Marshal(data)
	}
	if err != nil {
		return err
	}

	rsp, err := http.Post(notifier.URL, "application/json", bytes.NewBuffer(msg))
	if err != nil {
		return err
	}
	if len(rspHookFn) > 0 && rspHookFn[0] != nil {
		body, err = io.ReadAll(rsp.Body)
		if err != nil {
			return err
		}
		defer rsp.Body.Close()
		return rspHookFn[0](body)
	}
	return nil
}
