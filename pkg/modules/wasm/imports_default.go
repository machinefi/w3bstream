package wasm

import (
	"context"
	"os"

	"github.com/pkg/errors"
	"golang.org/x/exp/slog"

	"github.com/machinefi/w3bstream/pkg/modules/wasm/consts"
)

var ErrNotImplements = errors.New("not implements")

var (
	_ ImportsSystem         = (*DefaultImportsSystem)(nil)
	_ ImportsResource       = (*DefaultImportsResource)(nil)
	_ ImportsKVStore        = (*DefaultImportsKVStore)(nil)
	_ ImportsSQL            = (*DefaultImportsSQL)(nil)
	_ ImportsChainOperation = (*DefaultImportsChainOperation)(nil)
	_ ImportsMQTT           = (*DefaultImportsMQTT)(nil)
	_ ImportsMetrics        = (*DefaultImportsMetrics)(nil)
	_ ImportsAsyncCall      = (*DefaultImportsAsyncCall)(nil)
)

var DefaultImportsHandler = struct {
	ImportsSystem
	ImportsResource
	ImportsKVStore
	ImportsSQL
	ImportsChainOperation
	ImportsMQTT
	ImportsMetrics
	ImportsAsyncCall
}{
	ImportsSystem:         &DefaultImportsSystem{},
	ImportsResource:       &DefaultImportsResource{},
	ImportsKVStore:        &DefaultImportsKVStore{},
	ImportsSQL:            &DefaultImportsSQL{},
	ImportsChainOperation: &DefaultImportsChainOperation{},
	ImportsMQTT:           &DefaultImportsMQTT{},
	ImportsMetrics:        &DefaultImportsMetrics{},
	ImportsAsyncCall:      &DefaultImportsAsyncCall{},
}

type DefaultImportsSystem struct{}

func (l *DefaultImportsSystem) Log(level consts.LogLevel, msg string) {
	slog.Log(context.Background(), level.Level(), msg)
}

func (l *DefaultImportsSystem) LogInternal(level consts.LogLevel, msg string) {
	slog.Log(context.Background(), level.Level(), msg)
}

func (l *DefaultImportsSystem) Env(key string) (string, bool) {
	return os.LookupEnv(key)
}

type DefaultImportsResource struct{}

func (v *DefaultImportsResource) GetResourceData(uint32) ([]byte, bool) { return nil, false }

func (v *DefaultImportsResource) SetResourceData(uint32, []byte) error { return ErrNotImplements }

type DefaultImportsKVStore struct{}

func (*DefaultImportsKVStore) GetKVData(string) ([]byte, error) { return nil, ErrNotImplements }

func (*DefaultImportsKVStore) SetKVData(string, []byte) error { return ErrNotImplements }

type DefaultImportsSQL struct{}

func (*DefaultImportsSQL) ExecSQL(string) error { return ErrNotImplements }

func (*DefaultImportsSQL) QuerySQL(string) ([]byte, error) { return nil, ErrNotImplements }

type DefaultImportsChainOperation struct{}

func (*DefaultImportsChainOperation) SendTX(chainID int32, data []byte) (string, error) {
	return "", ErrNotImplements
}

func (*DefaultImportsChainOperation) SendTXWithOperator(chainID int32, data []byte) (string, error) {
	return "", ErrNotImplements
}

func (*DefaultImportsChainOperation) CallContract(chainID int32, data []byte) ([]byte, error) {
	return nil, ErrNotImplements
}

type DefaultImportsMQTT struct{}

func (*DefaultImportsMQTT) PubMQTT(topic string, message []byte) error { return ErrNotImplements }

type DefaultImportsMetrics struct{}

func (*DefaultImportsMetrics) SubmitMetrics(data []byte) error { return ErrNotImplements }

type DefaultImportsAsyncCall struct{}

func (*DefaultImportsAsyncCall) AsyncAPICall(req []byte) ([]byte, error) {
	return nil, ErrNotImplements
}
