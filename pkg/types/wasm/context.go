package wasm

type ContextHandler interface {
	Name() string
	GetImports() ImportsHandler
	SetImports(ImportsHandler)
	GetExports() ExportsHandler
	GetInstance() Instance
	SetInstance(Instance)
}

type ABI struct {
	Imports  ImportsHandler
	Instance Instance
}

func (a *ABI) Name() string { return NameVersion }

func (a *ABI) GetExports() ExportsHandler { return a }

func (a *ABI) GetImports() ImportsHandler { return a.Imports }

func (a *ABI) SetImports(i ImportsHandler) { a.Imports = i }

func (a *ABI) GetInstance() Instance { return a.Instance }

func (a *ABI) SetInstance(i Instance) { a.Instance = i }

func (a *ABI) Start() {}

func (a *ABI) Alloc() {}

func (a *ABI) Free() {}
