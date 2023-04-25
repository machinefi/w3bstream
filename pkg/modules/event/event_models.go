package event

import (
	"github.com/machinefi/w3bstream/pkg/depends/protocol/eventpb"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

type HandleEventResult struct {
	ProjectName string                    `json:"projectName"`
	PubID       types.SFID                `json:"pubID,omitempty"`
	PubName     string                    `json:"pubName,omitempty"`
	EventID     string                    `json:"eventID"`
	ErrMsg      string                    `json:"errMsg,omitempty"`
	WasmResults []*wasm.EventHandleResult `json:"wasmResults"`
}

type HandleEventReq struct {
	Events []eventpb.Event `json:"events"`
}
