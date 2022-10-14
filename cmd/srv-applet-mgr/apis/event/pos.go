package event

import (
	"context"
	"unicode/utf8"

	"github.com/iotexproject/Bumblebee/kit/httptransport/httpx"

	"github.com/iotexproject/w3bstream/pkg/depends/protocol/eventpb"
	"github.com/iotexproject/w3bstream/pkg/depends/unit"
	"github.com/iotexproject/w3bstream/pkg/enums"
	"github.com/iotexproject/w3bstream/pkg/errors/status"
	"github.com/iotexproject/w3bstream/pkg/modules/strategy"
	"github.com/iotexproject/w3bstream/pkg/modules/vm"
	"github.com/iotexproject/w3bstream/pkg/types"
	"github.com/iotexproject/w3bstream/pkg/types/wasm"
)

const (
	strLenLimit   = 50
	dataSizeLimit = 2 * unit.KiB
)

type HandleEvent struct {
	httpx.MethodPost
	ProjectName   string `in:"path" name:"projectName"`
	eventpb.Event `in:"body"`
}

func (r *HandleEvent) Path() string { return "/:projectName" }

func (r *HandleEvent) Output(ctx context.Context) (interface{}, error) {
	// TODO validate publisher belongs to Project @ZhiweiSun

	// TODO check publisher key
	//eventProto := models.EventProto{}
	//if err := json.Unmarshal(r.Data, &eventProto); err != nil {
	//
	//}
	//eventProto.PublisherKey -> projectID -> == r.ProjectID -> next

	eventType := enums.EVENT_TYPE__ANY
	if r.Header != nil {
		eventType = enums.EventType(r.Header.EventType)
	}

	instances, err := strategy.FindStrategyInstances(ctx, r.ProjectName, eventType)
	if err != nil {
		return nil, err
	}

	if len(r.Payload) > dataSizeLimit {
		return nil, status.BadRequest
	}

	ret := make([]HandleEventRsp, 0, len(instances))

	for _, v := range instances {
		consumer := vm.GetConsumer(v.InstanceID.String())
		if consumer == nil {
			continue
		}
		// TODO
		_, code := consumer.HandleEvent(v.Handler, []byte(r.Payload))
		ret = append(ret, HandleEventRsp{
			InstanceID: v.InstanceID,
			ResultCode: code,
		})
	}
	return ret, nil
}

type HandleEventRsp struct {
	InstanceID types.SFID            `json:"instanceID"`
	ResultCode wasm.ResultStatusCode `json:"resultCode"`
}

func check(projectID, appletID, publisher types.SFID, handler string) bool {
	return utf8.RuneCountInString(handler) <= strLenLimit
}
