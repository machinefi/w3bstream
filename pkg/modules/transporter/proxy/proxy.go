package proxy

import (
	"context"
	"strings"

	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/depends/protocol/eventpb"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/modules/event"
	"github.com/machinefi/w3bstream/pkg/types"
)

func Forward(ctx context.Context, channel string, ev *eventpb.Event) (interface{}, error) {
	cli := types.MustProxyClientFromContext(ctx)
	req := &event.EventReq{
		Channel:   channel,
		EventType: ev.Header.GetEventType(),
		EventID:   ev.Header.EventId,
		Payload:   ev.Payload,
	}

	tok := ev.Header.GetToken()
	if tok == "" {
		return nil, status.InvalidEventToken
	}

	if !strings.HasPrefix(tok, "Bearer") {
		tok = "Bearer " + tok
	}

	meta := kit.Metadata{
		"Authorization": []string{tok},
	}

	rsp := &event.EventRsp{}
	if _, err := cli.Do(ctx, req, meta).Into(rsp); err != nil {
		return nil, err
	}
	return rsp, nil
}
