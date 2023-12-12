package runtime

import "github.com/pkg/errors"

var (
	ErrInvalidWasmCode   = errors.New("")
	ErrNewWasmModule     = errors.New("")
	ErrInstanceNotStart  = errors.New("")
	ErrInvalidExport     = errors.New("")
	ErrInvalidMemory     = errors.New("")
	ErrInvalidFunction   = errors.New("")
	ErrMemAccessOverflow = errors.New("")
)
