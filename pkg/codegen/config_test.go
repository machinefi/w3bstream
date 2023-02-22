package codegen_test

import (
	"fmt"

	"github.com/machinefi/w3bstream/pkg/codegen"
	g "github.com/machinefi/w3bstream/pkg/depends/gen/codegen"
)

var (
	f *g.File
	c *codegen.Config
)

func init() {
	f = g.NewFile("example", "mock.go")
	c = &codegen.Config{}
}

func ExampleConfig_SnippetMain() {
	fmt.Println(string(c.SnippetMain(f).Bytes()))
	// Output:
	// func main() {
	// var (
	// ctx context.Context
	// l log.Logger
	// d sqlx.DBExecutor
	// ins wasmtime.Instance
	// err error
	// )
	// }
}

func ExampleConfig_SnippetFilterFunc() {
	fmt.Println(string(c.SnippetFilterFunc(f).Bytes()))
	// Output:
	// func FilterFunc(input any) bool {
	// var res bool
	// src, err := json.Marshal(input)
	// if err != nil {
	// l.Error(err)
	// }
	// code := ins.HandleEvent(ctx, start, src)
	// if code < 0 {
	// return res
	// }
	// ret, ok := ins.GetResource.uint32(code)()
	// if ok {
	// defer ins.RmvResource.ctx.uint32(code)()
	// }
	// if !ok {
	// return res
	// }
	// switch strings.ToLower(ret) {
	// case "true":
	// res = true
	// case "false":
	// res = false
	// default:
	// l.Warn(errors.New("the value does not support"))
	// }
	// return res
	// }
}

func ExampleConfig_SnippetMapFunc() {
	fmt.Println(string(c.SnippetMapFunc(f).Bytes()))
	// Output:
	//func MapFunc(ctx context.Context, input interface{}) (interface{}, error) {
	//var (
	//res models.Customer
	//ret []byte
	//ok bool
	//)
	//src, err := json.Marshal(input)
	//if err != nil {
	//l.Error(err)
	//}
	//code := ins.HandleEvent(ctx, start, src)
	//if code < 0 {
	//return nil, errors.New("the value does not support")
	//}
	//ret, ok = ins.GetResource(uint32(code))
	//if ok {
	//defer ins.RmvResource(ctx, uint32(code))
	//}
	//if !ok {
	//l.Error(err)
	//return nil, errors.New("the value does not support")
	//}
	//err = json.Unmarshal(ret, res)
	//return res, err
	//}
}

func ExampleConfig_SnippetGroupByKey() {
	fmt.Println(string(c.SnippetGroupByKey(f).Bytes()))
	// Output:
	//func GroupByKey(item rxgo.Item) string {
	//var (
	//res string
	//ret []byte
	//ok bool
	//)
	//src, err := json.Marshal(input)
	//if err != nil {
	//l.Error(err)
	//}
	//code := ins.HandleEvent(ctx, start, src)
	//if code < 0 {
	//return error
	//}
	//ret, ok = ins.GetResource(uint32(code))
	//if ok {
	//defer ins.RmvResource(ctx, uint32(code))
	//}
	//if !ok {
	//l.Error(err)
	//return error
	//}
	//res = string(ret)
	//return res
	//}
}
