package enums

//go:generate toolkit gen enum WasmDBDatatype
type WasmDBDatatype uint8

const (
	WASM_DB_DATATYPE_UNKNOWN WasmDBDatatype = iota
	WASM_DB_DATATYPE__INT
	WASM_DB_DATATYPE__INT8
	WASM_DB_DATATYPE__INT16
	WASM_DB_DATATYPE__INT32
	WASM_DB_DATATYPE__INT64
	WASM_DB_DATATYPE__UINT
	WASM_DB_DATATYPE__UINT8
	WASM_DB_DATATYPE__UINT16
	WASM_DB_DATATYPE__UINT32
	WASM_DB_DATATYPE__UINT64
	WASM_DB_DATATYPE__FLOAT32
	WASM_DB_DATATYPE__FLOAT64
	WASM_DB_DATATYPE__TEXT
	WASM_DB_DATATYPE__BOOL
	WASM_DB_DATATYPE__TIMESTAMP // use epoch timestamp (integer, UTC)
	WASM_DB_DATATYPE__DECIMAL
)
