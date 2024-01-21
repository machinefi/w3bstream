package wasm

import "github.com/machinefi/w3bstream/pkg/modules/wasm/consts"

type ImportsHandler interface {
	ImportsSystem
	ImportsResource
	ImportsKVStore
	ImportsSQL
	ImportsChainOperation
	ImportsMQTT
	ImportsMetrics
	ImportsAsyncCall
}

type ImportsSystem interface {
	// Log host log
	Log(level consts.LogLevel, msg string)
	// LogInternal host functions log
	LogInternal(level consts.LogLevel, msg string)
	// Env get env var by key
	Env(key string) (string, bool)
}

type ImportsResource interface {
	// GetResourceData get resource data by resource id
	GetResourceData(rid uint32) ([]byte, bool)
	// SetResourceData set resource data by resource id
	SetResourceData(rid uint32, data []byte) error
}

type ImportsKVStore interface {
	// GetKVData read value by key from kv store
	GetKVData(key string) ([]byte, error)
	// SetKVData set value by key to kv store
	SetKVData(key string, data []byte) error
}

type ImportsSQL interface {
	// ExecSQL exec sql query
	ExecSQL(q string) error
	// QuerySQL query data
	QuerySQL(q string) ([]byte, error)
}

type ImportsChainOperation interface {
	// SendTX send tx by chain id and data, returns tx hash and error
	SendTX(chainID int32, data []byte) (string, error)
	// SendTXWithOperator
	SendTXWithOperator(chainID int32, data []byte) (string, error)
	// CallContract call contract by chain id and data, returns call result
	CallContract(chainID int32, data []byte) ([]byte, error)
}

type ImportsMQTT interface {
	// PubMQTT publish mqtt message
	PubMQTT(topic string, message []byte) error
}

type ImportsMetrics interface {
	// SubmitMetrics submit metrics
	SubmitMetrics(data []byte) error
}

type ImportsAsyncCall interface {
	// AsyncAPICall call async api
	AsyncAPICall(req []byte) ([]byte, error)
}
