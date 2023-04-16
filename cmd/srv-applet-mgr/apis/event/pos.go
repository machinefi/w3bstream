package event

import (
	"context"
	"strings"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/depends/protocol/eventpb"
	"github.com/machinefi/w3bstream/pkg/modules/event"
	"github.com/machinefi/w3bstream/pkg/types"
)

type BatchEventHandleProxy struct {
	httpx.MethodPost
	Channel              string `in:"path" name:"channel"`
	event.HandleEventReq `in:"body"`
}

func (r *BatchEventHandleProxy) Path() string { return "/proxy/:channel" }

func (r *BatchEventHandleProxy) Output(ctx context.Context) (interface{}, error) {
	results := make([]*event.HandleEventResult, len(r.Events))
	for i := range r.Events {
		evt := &r.Events[i]
		res, err := ProxyForward(ctx, r.Channel, evt)
		if err != nil {
			res = &event.HandleEventResult{
				ProjectName: r.Channel,
				EventID:     evt.Header.GetEventId(),
				ErrMsg:      err.Error(),
			}
		}
		results = append(results, res)
	}
	return results, nil
}

func ProxyForward(ctx context.Context, channel string, ev *eventpb.Event) (*event.HandleEventResult, error) {
	cli := types.MustEventProxyClientFromContext(ctx)
	req := &EventHandle{
		Channel:   channel,
		EventType: ev.Header.GetEventType(),
		EventID:   ev.Header.EventId,
		Payload:   ev.Payload,
	}
	meta := kit.Metadata{}
	tok := ev.Header.GetToken()
	if tok != "" {
		if !strings.HasPrefix(tok, "Bearer") {
			tok = "Bearer " + tok
		}
		meta.Add("Authorization", tok)
	}

	rsp := &event.HandleEventResult{}
	if _, err := cli.Do(ctx, req, meta).Into(rsp); err != nil {
		return nil, err
	}
	return rsp, nil
}

type EventHandle struct {
	httpx.MethodPost
	Channel   string `in:"path" name:"channel"`
	EventType string `in:"path" name:"eventType"`
	EventID   string `in:"path" name:"eventID"`
	Payload   []byte `in:"body"`
}

func (r *EventHandle) Path() string {
	return "/transport/:channel/:eventType/:eventID"
}

func (r *EventHandle) Output(ctx context.Context) (interface{}, error) {
	var (
		err error
		pub = middleware.MustPublisher(ctx)
		rsp = &event.HandleEventResult{
			ProjectName: r.Channel,
			PubID:       pub.PublisherID,
			PubName:     pub.Name,
			EventID:     r.EventID,
		}
	)

	ctx, err = pub.WithStrategiesContextByChannelAndType(ctx, r.Channel, r.EventType)
	if err != nil {
		rsp.ErrMsg = err.Error()
		return rsp, nil
	}
	rsp.WasmResults = event.OnEventReceived(ctx, r.Payload)
	return rsp, nil
}
