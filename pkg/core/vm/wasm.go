package vm

/*
import w "github.com/wasmerio/wasmer-go/wasmer"

type Wasm struct {
	code     []byte
	engine   *w.Engine
	store    *w.Store
	module   *w.Module
	instance *w.Instance
}

type NativeFunc = func(args []w.Value) ([]w.Value, error)

type WasmImportFunc struct {
	name        string
	inputTypes  []w.ValueKind
	outputTypes []w.ValueKind
	nativeFunc  NativeFunc
}

type WasmImport struct {
	namespace string
	functions []WasmImportFunc
}

func NewWasm(code []byte, imports []WasmImport) (*Wasm, error) {
	var instance *w.Instance

	engine := w.NewEngine()
	store := w.NewStore(engine)
	module, _ := w.NewModule(store, code)

	importObject := w.NewImportObject()

	for _, v := range imports {
		intoExtern := make(map[string]w.IntoExtern)
		for _, fn := range v.functions {
			intoExtern[fn.name] = w.NewFunction(
				store,
				w.NewFunctionType(w.NewValueTypes(fn.inputTypes...), w.NewValueTypes(fn.outputTypes...)),
				fn.nativeFunc,
			)
		}
		importObject.Register(
			v.namespace,
			intoExtern,
		)
	}

	instance, e := w.NewInstance(module, importObject)
	if e != nil {
		return nil, e
	}

	return &Wasm{
		code:     code,
		engine:   engine,
		store:    store,
		module:   module,
		instance: instance,
	}, nil
}

func (wasm *Wasm) GetFunction(name string) (w.NativeFunction, error) {
	return wasm.instance.Exports.GetFunction(name)
}

func (wasm *Wasm) ExecuteFunction(name string, args ...interface{}) (interface{}, error) {
	fn, e := wasm.instance.Exports.GetFunction(name)
	if e != nil {
		return nil, e
	}
	return fn(args...)
}

func (wasm *Wasm) GetMemory(name string) ([]byte, error) {
	memory, e := wasm.instance.Exports.GetMemory(name)
	if e != nil {
		return nil, e
	}
	return memory.Data(), nil
}
*/
