package wasm

import "github.com/machinefi/w3bstream/pkg/enums"

var NameVersion = "w3bstream@v0.0.1"

// ResultStatusCode wasm call result code
type ResultStatusCode int32

const (
	ResultStatusCode_OK ResultStatusCode = iota
	ResultStatusCode_UnexportedHandler
	ResultStatusCode_ResourceNotFound
	ResultStatusCode_ImportNotFound
	ResultStatusCode_ImportCallFailed
	ResultStatusCode_TransDataToVMFailed
	ResultStatusCode_TransDataFromVMFailed
	ResultStatusCode_HostInternal

	// TODO following result status
	ResultStatusCode_Failed = -1 // reserved for wasm invoke failed
)

type InstanceState = enums.InstanceState
