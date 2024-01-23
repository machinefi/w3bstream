package wasm

type Context interface {
	Name() string

	GetImports() ImportsHandler
	SetImports(ImportsHandler)

	GetExports() Exports

	GetInstance() Instance
	SetInstance(Instance)
}
