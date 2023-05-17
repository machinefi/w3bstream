package proxy

import (
	"bytes"
	"context"
	"strings"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/depends/protocol/eventpb"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/modules/event"
	"github.com/machinefi/w3bstream/pkg/types"
)

type ForwardRequest struct {
	httpx.MethodPost
	event.EventReq
}

func (r *ForwardRequest) Path() string {
	return "/srv-applet-mgr/v0/event/" + r.Channel
}

func Forward(ctx context.Context, channel string, ev *eventpb.Event) (interface{}, error) {
	cli := types.MustProxyClientFromContext(ctx)

	body := event.EventReq{
		Channel:   channel,
		EventType: ev.Header.GetEventType(),
		EventID:   ev.Header.GetEventId(),
		Payload:   *(bytes.NewBuffer(ev.Payload)),
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
	req := &ForwardRequest{EventReq: body}
	if _, err := cli.Do(context.Background(), req, meta).Into(rsp); err != nil {
		return nil, err
	}
	return rsp, nil
}
