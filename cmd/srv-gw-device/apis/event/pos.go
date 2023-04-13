package event

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/machinefi/w3bstream/cmd/srv-gw-device/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/device"
	"github.com/machinefi/w3bstream/pkg/modules/event"
)

type HandleEvent struct {
	httpx.MethodPost
	event.EventGateway `in:"body"`
}

func (r *HandleEvent) Path() string { return "/" }

func (r *HandleEvent) Output(ctx context.Context) (interface{}, error) {
	d, err := device.GetDeviceByID(ctx, r.DeviceID)
	if err != nil {
		return nil, err
	}
	// can do more device check here
	return r.callCore(ctx, d)
}

func (r *HandleEvent) callCore(ctx context.Context, d *models.Device) (interface{}, error) {
	core := types.MustCoreFromContext(ctx)
	l := types.MustLoggerFromContext(ctx)

	_, l = l.Start(ctx, "GatewayCallCore")
	defer l.End()

	e := &event.EventCore{
		ProjectID: d.ProjectID,
		EventID:   r.EventID,
		EventType: r.EventType,
		Payload:   r.Payload,
	}
	body, err := json.Marshal(e)
	if err != nil {
		l.Error(err)
		return nil, err
	}

	req, err := http.NewRequest("POST", core.Endpoint.String(), bytes.NewBuffer(body))
	if err != nil {
		l.Error(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		l.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		err := fmt.Errorf("unexpected http code %v", resp.StatusCode)
		l.Error(err)
		return nil, err
	}
	return resp, nil
}
