package wasmtime

import "github.com/pkg/errors"

var (
	ErrInvalidWasmCode       = errors.New("invalid wasm code")
	ErrFailedToNewWasmModule = errors.New("failed to new wasm module")
	ErrInstanceNotStarted    = errors.New("instance not started")
	ErrInvalidExportFunc     = errors.New("invalid export func")
	ErrInvalidExportMem      = errors.New("invalid export memory")
	ErrInvalidImportFunc     = errors.New("invalid import func")
	ErrMemAccessOverflow     = errors.New("memory access overflow")
	ErrUnknownABIName        = errors.New("unknown abi name")
)
