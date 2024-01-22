package consts

//go:generate toolkit gen enum Result
type Result int32

const (
	RESULT_UNKNOWN Result = iota
	RESULT__INVALID_MEM_ACCESS
	RESULT__ENV_NOT_FOUND
	RESULT__RESOURCE_NOT_FOUND
	RESULT__RESOURCE_EVENT_NOT_FOUND
	RESULT__KV_DATA_NOT_FOUND
	RESULT__IMPORT_HANDLE_FAILED
	RESULT__HOST_INVOKE_FAILED
)

const RESULT_OK = RESULT_UNKNOWN

func (v Result) OK() bool {
	return v == RESULT_OK
}

func (v Result) Error() string {
	return v.String()
}

func (v Result) Int32() int32 {
	return int32(v)
}
