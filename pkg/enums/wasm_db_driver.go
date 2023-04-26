package enums

//go:generate toolkit gen enum WasmDBDialect
type WasmDBDialect uint8

const (
	WASM_DB_DIALECT_UNKNOWN WasmDBDialect = iota
	WASM_DB_DIALECT__POSTGRES
)
