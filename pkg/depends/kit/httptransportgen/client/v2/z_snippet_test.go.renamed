package client_test

import (
	"fmt"

	g "github.com/machinefi/w3bstream/pkg/depends/gen/codegen"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransportgen/client"
)

func output(ss ...g.Snippet) {
	for _, s := range ss {
		fmt.Println(string(s.Bytes()))
	}
	fmt.Println()
}

var f = g.NewFile("mock", "mock")

func ExampleSnippetOperationDefine() {
	s := client.SnippetOperationDefine(
		"Cookie",
		&g.SnippetField{
			Type:  g.String,
			Names: []*g.SnippetIdent{g.Ident("Token")},
			Tag:   `in:"cookie" name:"token,omitempty"`,
		},
	)
	output(s)
	// Output:
	// type Cookie struct {
	// Token string `in:"cookie" name:"token,omitempty"`
	// }
}

func ExampleSnippetOperationPathMethod() {
	s := client.SnippetOperationPathMethod(f, "Cookie", "/demo/cookie")
	output(s)
	// Output:
	// func (o *Cookie) Path() string {
	// return "/demo/cookie"
	// }
}

func ExampleSnippetOperationMethodMethod() {
	s := client.SnippetOperationMethodMethod(f, "Cookie", "POST")
	output(s)
	// Output:
	// func (o *Cookie) Method() string {
	// return "POST"
	// }
}

func ExampleSnippetOperationDoMethod() {
	s := client.SnippetOperationDoMethod(f, "demo", "Cookie", "comment 1", "comment 2")
	output(s...)
	// Output:
	// // comment 1
	// // comment 2
	// func (o *Cookie) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	// ctx = metax.ContextWith(ctx, "operationID", "demo.Cookie")
	// return cli.Do(ctx, o, metas...)
	// }
}

func ExampleSnippetOperationInvokeContextMethod() {
	output(
		client.SnippetOperationInvokeContextMethod(f, "Cookie", g.Type("FAKE_RESPONSE")),
		client.SnippetOperationInvokeContextMethod(f, "Cookie", nil),
	)
	// Output:
	// func (o *Cookie) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (*FAKE_RESPONSE, kit.Metadata, error) {
	// rsp := new(FAKE_RESPONSE)
	// meta, err := cli.Do(ctx, o, metas...).Into(rsp)
	// return rsp, meta, err
	// }
	// func (o *Cookie) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	// meta, err := cli.Do(ctx, o, metas...).Into(nil)
	// return meta, err
	// }
}

func ExampleSnippetOperationInvokeMethod() {
	output(
		client.SnippetOperationInvokeMethod(f, "Cookie", g.Type("FAKE_RESPONSE")),
		client.SnippetOperationInvokeMethod(f, "Cookie", nil),
	)
	// Output:
	// func (o *Cookie) Invoke(cli kit.Client, metas ...kit.Metadata) (*FAKE_RESPONSE, kit.Metadata, error) {
	// return o.InvokeContext(context.Background(), cli, metas...)
	// }
	// func (o *Cookie) Invoke(cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	// return o.InvokeContext(context.Background(), cli, metas...)
	// }
}

func ExampleSnippetClientInterface() {
	output(
		client.SnippetClientInterface(
			f,
			g.Func().Named("ABC"), g.Func().Named("DEF"),
		),
	)
	// Output:
	// type Interface interface {
	// Context() context.Context
	// WithContext(context.Context) Interface
	// ABC()
	// DEF()
	// }
}

func ExampleSnippetNewClient() {
	output(client.SnippetNewClient(f))
	// Output:
	// func NewClient(c kit.Client) *Client {
	// return &(Client{
	// Client: c,
	// })
	// }
}

func ExampleSnippetClientDefine() {
	output(client.SnippetClientDefine(f))
	// Output:
	// type Client struct {
	// Client kit.Client
	// ctx context.Context
	// }
}

func ExampleSnippetClientContextMethod() {
	output(client.SnippetClientContextMethod(f))
	// Output:
	// func (c *Client) Context() context.Context {
	// if c.ctx != nil {
	// return c.ctx
	// }
	// return context.Background()
	// }
}

func ExampleSnippetClientWithContextMethod() {
	output(client.SnippetClientWithContextMethod(f))
	// Output:
	// func (c *Client) WithContext(ctx context.Context) Interface {
	// cc := new(Client)
	// cc.Client, cc.ctx = c.Client, ctx
	// return cc
	// }
}
