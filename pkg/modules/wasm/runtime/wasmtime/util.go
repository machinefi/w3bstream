package wasmtime

import (
	"reflect"

	"github.com/bytecodealliance/wasmtime-go/v8"
	"github.com/pkg/errors"
)

func checkIfOverflow(addr, size int32, mem []byte) bool {
	return addr > int32(len(mem)) || addr+size > int32(len(mem))
}

func newWasmtimeNativeFunction(store *wasmtime.Store, f *wasmtime.Func) *wasmtimeNativeFunction {
	return &wasmtimeNativeFunction{store: store, f: f}
}

type wasmtimeNativeFunction struct {
	store *wasmtime.Store
	f     *wasmtime.Func
}

func (wf *wasmtimeNativeFunction) Call(args ...any) (any, error) {
	return wf.f.Call(wf.store, args...)
}

func (wf *wasmtimeNativeFunction) Params() []string {
	params := wf.f.Type(wf.store).Params()

	ret := make([]string, 0, len(params))
	for _, t := range params {
		ret = append(ret, t.String())
	}
	return ret
}

func (wf *wasmtimeNativeFunction) Results() []string {
	results := wf.f.Type(wf.store).Results()

	ret := make([]string, 0, len(results))
	for _, t := range results {
		ret = append(ret, t.String())
	}
	return ret
}

func convertFromGoType(t reflect.Type) *wasmtime.ValType {
	switch t.Kind() {
	case reflect.Int32:
		return wasmtime.NewValType(wasmtime.KindI32)
	case reflect.Int64:
		return wasmtime.NewValType(wasmtime.KindI64)
	case reflect.Float32:
		return wasmtime.NewValType(wasmtime.KindF32)
	case reflect.Float64:
		return wasmtime.NewValType(wasmtime.KindF64)
	default:
		panic(errors.Errorf("convertFromGoType unsupported go type: %s", t))
	}
}

func convertToGoTypes(in wasmtime.Val) reflect.Value {
	switch in.Kind() {
	case wasmtime.KindI32:
		return reflect.ValueOf(in.I32())
	case wasmtime.KindI64:
		return reflect.ValueOf(in.I64())
	case wasmtime.KindF32:
		return reflect.ValueOf(in.F32())
	case wasmtime.KindF64:
		return reflect.ValueOf(in.F64())
	default:
		panic(errors.Errorf("convertToGoType unsupported go type: %s", in.Kind().String()))
	}
}

func convertToWasmtimeVal(v interface{}) wasmtime.Val {
	switch _v := v.(type) {
	case int32:
		return wasmtime.ValI32(_v)
	case int64:
		return wasmtime.ValI64(_v)
	case float32:
		return wasmtime.ValF32(_v)
	case float64:
		return wasmtime.ValF64(_v)
	case reflect.Value:
		return convertToWasmtimeVal(_v.Interface())
	default:
		panic(errors.Errorf("convertToWasmtimeVal unsupported go type: %s", reflect.TypeOf(v)))
	}
}
