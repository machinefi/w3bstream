package wasm

type ExportsHandler interface {
	Start()
	Alloc()
	Free()
}
