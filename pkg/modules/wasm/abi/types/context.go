package types

type Context interface {
	Name() string

	GetImports() ImportsHandler
	SetImports(ImportsHandler)

	GetExports() Exports

	GetInstance() Instance
	SetInstance(Instance)

	OnCreated(Instance)
	OnDestroy(Instance)
}
