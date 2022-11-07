package wasmtime_test

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	// . "github.com/onsi/gomega"
)

var (
	wasmLogCode             []byte
	wasmGJsonCode           []byte
	wasmEasyJsonCode        []byte
	wasmWordCountCode       []byte
	wasmWordCountV2Code     []byte
	wasmTokenDistributeCode []byte

	// ctx context.Context
	// idg confid.SFIDGenerator
)

func init() {
	wd, _ := os.Getwd()
	fmt.Println(wd)
	root := filepath.Join(wd, "../../../../_examples")
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

	/*
		ctx = global.WithContext(context.Background())
		ctx = types.WithETHClientConfig(ctx, &types.ETHClientConfig{
			PrivateKey:    "",
			ChainEndpoint: "https://babel-api.testnet.iotex.io",
		})

		idg = confid.MustSFIDGeneratorFromContext(ctx)

		go kit.Run(tasks.Root, global.TaskServer())
	*/
}

/*
func TestInstance_LogWASM(t *testing.T) {
	i, err := wasmtime.NewInstanceByCode(ctx, idg.MustGenSFID(), wasmLogCode)
	NewWithT(t).Expect(err).To(BeNil())
	id := vm.AddInstance(ctx, i)

	err = vm.StartInstance(ctx, id)
	NewWithT(t).Expect(err).To(BeNil())
	defer vm.StopInstance(ctx, id)

	ret := i.HandleEvent(ctx, "start", []byte("IoTeX"))
	NewWithT(t).Expect(ret.Code).To(Equal(wasm.ResultStatusCode_OK))
	NewWithT(t).Expect(ret.ErrMsg).To(Equal(""))

	ret = i.HandleEvent(ctx, "not_exported", []byte("IoTeX"))
	NewWithT(t).Expect(ret.Code).To(Equal(wasm.ResultStatusCode_UnexportedHandler))
}

func TestInstance_GJsonWASM(t *testing.T) {
	i, err := wasmtime.NewInstanceByCode(ctx, idg.MustGenSFID(), wasmGJsonCode)
	NewWithT(t).Expect(err).To(BeNil())
	id := vm.AddInstance(ctx, i)

	err = vm.StartInstance(ctx, id)
	NewWithT(t).Expect(err).To(BeNil())
	defer vm.StopInstance(ctx, id)

	ret := i.HandleEvent(ctx, "start", []byte(`
{
  "name": {
    "first": "Tom",
    "last": "Anderson",
    "age": 39
  },
  "friends": [
    {
      "first_name": "Dale",
      "last_name": "Murphy",
      "age": 44,
      "nets": [
        "ig",
        "fb",
        "tw"
      ]
    },
    {
      "first_name": "Roger",
      "last_name": "Craig",
      "age": 68,
      "nets": [
        "fb",
        "tw"
      ]
    },
    {
      "first_name": "Jane",
      "last_name": "Murphy",
      "age": 47,
      "nets": [
        "ig",
        "tw"
      ]
    }
  ]
}`))
	NewWithT(t).Expect(ret.Code).To(Equal(wasm.ResultStatusCode_OK))
}

func TestInstance_EasyJsonWASM(t *testing.T) {
	i, err := wasmtime.NewInstanceByCode(ctx, idg.MustGenSFID(), wasmEasyJsonCode)
	NewWithT(t).Expect(err).To(BeNil())
	id := vm.AddInstance(ctx, i)

	err = vm.StartInstance(ctx, id)
	NewWithT(t).Expect(err).To(BeNil())
	defer vm.StopInstance(ctx, id)

	ret := i.HandleEvent(ctx, "start", []byte(`
{
  "id": 11,
  "student_name": "Tom",
  "student_school": {
    "school_name": "MIT",
    "school_addr": "xz"
  },
  "birthday": "2017-08-04T20:58:07.9894603+08:00"
}`))
	NewWithT(t).Expect(ret.Code).To(Equal(wasm.ResultStatusCode_OK))
}

func TestInstance_WordCount(t *testing.T) {
	i, err := wasmtime.NewInstanceByCode(ctx, idg.MustGenSFID(), wasmWordCountCode)
	NewWithT(t).Expect(err).To(BeNil())
	id := vm.AddInstance(ctx, i)

	err = vm.StartInstance(ctx, id)
	NewWithT(t).Expect(err).To(BeNil())
	defer vm.StopInstance(ctx, id)

	ret := i.HandleEvent(ctx, "start", []byte("a b c d a"))
	NewWithT(t).Expect(ret.Code).To(Equal(wasm.ResultStatusCode_OK))

	NewWithT(t).Expect(i.Get("a")).To(Equal(int32(2)))
	NewWithT(t).Expect(i.Get("b")).To(Equal(int32(1)))
	NewWithT(t).Expect(i.Get("c")).To(Equal(int32(1)))
	NewWithT(t).Expect(i.Get("d")).To(Equal(int32(1)))

	ret = i.HandleEvent(ctx, "start", []byte("a b c d a"))
	NewWithT(t).Expect(ret.Code).To(Equal(wasm.ResultStatusCode_OK))

	NewWithT(t).Expect(i.Get("a")).To(Equal(int32(4)))
	NewWithT(t).Expect(i.Get("b")).To(Equal(int32(2)))
	NewWithT(t).Expect(i.Get("c")).To(Equal(int32(2)))
	NewWithT(t).Expect(i.Get("d")).To(Equal(int32(2)))
}

func TestInstance_WordCountV2(t *testing.T) {
	i, err := wasmtime.NewInstanceByCode(ctx, idg.MustGenSFID(), wasmWordCountV2Code)
	NewWithT(t).Expect(err).To(BeNil())
	id := vm.AddInstance(ctx, i)

	err = vm.StartInstance(ctx, id)
	NewWithT(t).Expect(err).To(BeNil())
	defer vm.StopInstance(ctx, id)

	ret := i.HandleEvent(ctx, "start", []byte("a b c d a"))
	NewWithT(t).Expect(ret.Code).To(Equal(wasm.ResultStatusCode_OK))

	NewWithT(t).Expect(i.Get("a")).To(Equal(int32(2)))
	NewWithT(t).Expect(i.Get("b")).To(Equal(int32(1)))
	NewWithT(t).Expect(i.Get("c")).To(Equal(int32(1)))
	NewWithT(t).Expect(i.Get("d")).To(Equal(int32(1)))

	ret = i.HandleEvent(ctx, "start", []byte("a b c d a"))
	NewWithT(t).Expect(ret.Code).To(Equal(wasm.ResultStatusCode_OK))

	NewWithT(t).Expect(i.Get("a")).To(Equal(int32(4)))
	NewWithT(t).Expect(i.Get("b")).To(Equal(int32(2)))
	NewWithT(t).Expect(i.Get("c")).To(Equal(int32(2)))
	NewWithT(t).Expect(i.Get("d")).To(Equal(int32(2)))

	ret = i.HandleEvent(ctx, "word_count", nil)
	NewWithT(t).Expect(ret.Code).To(Equal(wasm.ResultStatusCode(4)))
}

func TestInstance_TokenDistribute(t *testing.T) {
	i, err := wasmtime.NewInstanceByCode(ctx, idg.MustGenSFID(), wasmTokenDistributeCode)
	NewWithT(t).Expect(err).To(BeNil())
	id := vm.AddInstance(ctx, i)

	err = vm.StartInstance(ctx, id)
	NewWithT(t).Expect(err).To(BeNil())
	defer vm.StopInstance(ctx, id)

	for idx := int32(0); idx < 20; idx++ {
		ret := i.HandleEvent(ctx, "start", []byte("test"))
		NewWithT(t).Expect(ret.Code).To(Equal(wasm.ResultStatusCode_OK))
		NewWithT(t).Expect(i.Get("clicks")).To(Equal(idx + 1))
	}
}
*/
