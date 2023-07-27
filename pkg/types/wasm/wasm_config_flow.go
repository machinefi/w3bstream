package wasm

import (
	"context"
	"encoding/json"

	"github.com/machinefi/w3bstream/pkg/enums"
)

type Flow struct {
	Source    Source     `json:"source"`
	Operators []Operator `json:"operators"`
	Sink      Sink       `json:"sink"`
}

func (f *Flow) ConfigType() enums.ConfigType {
	return enums.CONFIG_TYPE__PROJECT_FLOW
}

func (f *Flow) WithContext(ctx context.Context) context.Context {
	return WithFlow(ctx, f)
}

type Source struct {
	Strategies []string `json:"strategies"`
}

type Operator struct {
	OpType   enums.FlowOperator `json:"opType"`
	WasmFunc string             `json:"wasmFunc,omitempty"`
	Parallel int                `json:"parallel,omitempty"`
}

type Sink struct {
	SinkType enums.FlowSink `json:"sinkType"`
	SinkInfo SinkInfo       `json:"sinkInfo"`
}

type SinkInfo struct {
	DBInfo     DBInfo `json:"DBInfo,omitempty"`
	ChainBlock CBInfo `json:"chainBlock,omitempty"`
}

type DBInfo struct {
	Endpoint string   `json:"endpoint,omitempty"`
	DBType   string   `json:"DBType,omitempty"`
	Table    string   `json:"table,omitempty"`
	Columns  []string `json:"columns,omitempty"`
}

type CBInfo struct {
	ChainID int `json:"chainID,omitempty"`
}

func (s *Sink) UnmarshalJSON(b []byte) error {
	type Alias Sink
	var tmp Alias
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}

	s.SinkType = tmp.SinkType
	switch s.SinkType {
	case enums.FLOW_SINK__RMDB:
		s.SinkInfo = SinkInfo{DBInfo: tmp.SinkInfo.DBInfo}
	case enums.FLOW_SINK__BLOCKCHAIN:
		s.SinkInfo = SinkInfo{ChainBlock: tmp.SinkInfo.ChainBlock}
	case enums.FLOW_SINK_UNKNOWN:
		s.SinkInfo = tmp.SinkInfo
	default:
		panic("sink not support")
	}

	return nil
}
