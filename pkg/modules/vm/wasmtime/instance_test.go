package wasmtime_test

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/require"

	"github.com/iotexproject/w3bstream/pkg/modules/vm"
	"github.com/iotexproject/w3bstream/pkg/modules/vm/wasmtime"
	"github.com/iotexproject/w3bstream/pkg/types/wasm"
)

var (
	wasmLogCode             []byte
	wasmGJsonCode           []byte
	wasmEasyJsonCode        []byte
	wasmWordCountCode       []byte
	wasmWordCountV2Code     []byte
	wasmTokenDistributeCode []byte
)

func init() {
	wd, _ := os.Getwd()
	fmt.Println(wd)
	root := filepath.Join(wd, "../../../../examples")
	fmt.Println(root)

	var err error
	wasmLogCode, err = os.ReadFile(filepath.Join(root, "log/log.wasm"))
	if err != nil {
		panic(err)
	}

	wasmGJsonCode, err = os.ReadFile(filepath.Join(root, "gjson/gjson.wasm"))
	if err != nil {
		panic(err)
	}
	wasmEasyJsonCode, err = os.ReadFile(filepath.Join(root, "easyjson/easyjson.wasm"))
	if err != nil {
		panic(err)
	}
	wasmWordCountCode, err = os.ReadFile(filepath.Join(root, "word_count/word_count.wasm"))
	if err != nil {
		panic(err)
	}
	wasmWordCountV2Code, err = os.ReadFile(filepath.Join(root, "word_count_v2/word_count_v2.wasm"))
	if err != nil {
		panic(err)
	}

	wasmTokenDistributeCode, err = os.ReadFile(filepath.Join(root, "token_distribute/token_distribute.wasm"))
	if err != nil {
		panic(err)
	}
}

func TestInstance_LogWASM(t *testing.T) {
	require := require.New(t)
	i, err := wasmtime.NewInstanceByCode(context.Background(), wasmLogCode)
	require.NoError(err)
	id := vm.AddInstance(i)

	err = vm.StartInstance(id)
	require.NoError(err)
	defer vm.StopInstance(id)

	_, code := i.HandleEvent("start", []byte("IoTeX"))
	NewWithT(t).Expect(code).To(Equal(wasm.ResultStatusCode_OK))

	_, code = i.HandleEvent("not_exported", []byte("IoTeX"))
	NewWithT(t).Expect(code).To(Equal(wasm.ResultStatusCode_UnexportedHandler))
}

func TestInstance_GJsonWASM(t *testing.T) {
	require := require.New(t)
	i, err := wasmtime.NewInstanceByCode(context.Background(), wasmGJsonCode)
	require.NoError(err)
	id := vm.AddInstance(i)

	err = vm.StartInstance(id)
	require.NoError(err)
	defer vm.StopInstance(id)

	_, code := i.HandleEvent("start", []byte(`{
  "name": {"first": "Tom", "last": "Anderson", "age": 39},
  "friends": [
    {"first_name": "Dale", "last_name": "Murphy", "age": 44, "nets": ["ig", "fb", "tw"]},
    {"first_name": "Roger", "last_name": "Craig", "age": 68, "nets": ["fb", "tw"]},
    {"first_name": "Jane", "last_name": "Murphy", "age": 47, "nets": ["ig", "tw"]}
  ]
}`))
	NewWithT(t).Expect(code).To(Equal(wasm.ResultStatusCode_OK))
}

func TestInstance_EasyJsonWASM(t *testing.T) {
	require := require.New(t)
	i, err := wasmtime.NewInstanceByCode(context.Background(), wasmEasyJsonCode)
	require.NoError(err)
	id := vm.AddInstance(i)

	err = vm.StartInstance(id)
	require.NoError(err)
	defer vm.StopInstance(id)

	_, code := i.HandleEvent("start", []byte(`{"id":11,"student_name":"Tom","student_school":
								{"school_name":"MIT","school_addr":"xz"},
								"birthday":"2017-08-04T20:58:07.9894603+08:00"}`))
	NewWithT(t).Expect(code).To(Equal(wasm.ResultStatusCode_OK))
}

func TestInstance_WordCount(t *testing.T) {
	i, err := wasmtime.NewInstanceByCode(context.Background(), wasmWordCountCode)
	NewWithT(t).Expect(err).To(BeNil())
	id := vm.AddInstance(i)

	err = vm.StartInstance(id)
	NewWithT(t).Expect(err).To(BeNil())
	defer vm.StopInstance(id)

	_, code := i.HandleEvent("start", []byte("a b c d a"))
	NewWithT(t).Expect(code).To(Equal(wasm.ResultStatusCode_OK))

	NewWithT(t).Expect(i.Get("a")).To(Equal(int32(2)))
	NewWithT(t).Expect(i.Get("b")).To(Equal(int32(1)))
	NewWithT(t).Expect(i.Get("c")).To(Equal(int32(1)))
	NewWithT(t).Expect(i.Get("d")).To(Equal(int32(1)))

	_, code = i.HandleEvent("start", []byte("a b c d a"))
	NewWithT(t).Expect(code).To(Equal(wasm.ResultStatusCode_OK))

	NewWithT(t).Expect(i.Get("a")).To(Equal(int32(4)))
	NewWithT(t).Expect(i.Get("b")).To(Equal(int32(2)))
	NewWithT(t).Expect(i.Get("c")).To(Equal(int32(2)))
	NewWithT(t).Expect(i.Get("d")).To(Equal(int32(2)))
}

func TestInstance_WordCountV2(t *testing.T) {
	i, err := wasmtime.NewInstanceByCode(context.Background(), wasmWordCountV2Code)
	NewWithT(t).Expect(err).To(BeNil())
	id := vm.AddInstance(i)

	err = vm.StartInstance(id)
	NewWithT(t).Expect(err).To(BeNil())
	defer vm.StopInstance(id)

	_, code := i.HandleEvent("start", []byte("a b c d a"))
	NewWithT(t).Expect(code).To(Equal(wasm.ResultStatusCode_OK))

	NewWithT(t).Expect(i.Get("a")).To(Equal(int32(2)))
	NewWithT(t).Expect(i.Get("b")).To(Equal(int32(1)))
	NewWithT(t).Expect(i.Get("c")).To(Equal(int32(1)))
	NewWithT(t).Expect(i.Get("d")).To(Equal(int32(1)))

	_, code = i.HandleEvent("start", []byte("a b c d a"))
	NewWithT(t).Expect(code).To(Equal(wasm.ResultStatusCode_OK))

	NewWithT(t).Expect(i.Get("a")).To(Equal(int32(4)))
	NewWithT(t).Expect(i.Get("b")).To(Equal(int32(2)))
	NewWithT(t).Expect(i.Get("c")).To(Equal(int32(2)))
	NewWithT(t).Expect(i.Get("d")).To(Equal(int32(2)))

	_, unique := i.HandleEvent("word_count", nil)
	NewWithT(t).Expect(unique).To(Equal(wasm.ResultStatusCode(4)))
}

func TestInstance_TokenDistribute(t *testing.T) {
	i, err := wasmtime.NewInstanceByCode(context.Background(), wasmTokenDistributeCode)
	NewWithT(t).Expect(err).To(BeNil())
	id := vm.AddInstance(i)

	err = vm.StartInstance(id)
	NewWithT(t).Expect(err).To(BeNil())
	defer vm.StopInstance(id)

	for idx := int32(0); idx < 20; idx++ {
		_, code := i.HandleEvent("start", []byte("test"))
		NewWithT(t).Expect(code).To(Equal(wasm.ResultStatusCode_OK))
		NewWithT(t).Expect(i.Get("clicks")).To(Equal(idx + 1))
	}
}
