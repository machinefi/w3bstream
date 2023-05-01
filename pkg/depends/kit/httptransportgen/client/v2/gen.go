package client

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/machinefi/w3bstream/pkg/depends/gen/codegen"
	"github.com/machinefi/w3bstream/pkg/depends/oas"
	"github.com/machinefi/w3bstream/pkg/depends/x/stringsx"
	"github.com/pkg/errors"
)

type OasGenerator interface {
	Gen(ctx context.Context, spec *oas.OpenAPI) error
}

var (
	_ OasGenerator = (*ClientGen)(nil)
	_ OasGenerator = (*OperationGen)(nil)
	_ OasGenerator = (*TypeGen)(nil)
)

func NewGenerator(name string, u *url.URL, opts ...WithOption) *Generator {
	g := &Generator{
		ServiceName: name,
		URL:         u,
		Spec:        &oas.OpenAPI{},
	}

	for _, o := range opts {
		o(&g.Option)
	}

	return g
}

type Generator struct {
	ServiceName string
	URL         *url.URL
	Spec        *oas.OpenAPI

	Option
}

func (g *Generator) Load() {
	if g.URL == nil {
		panic(errors.Errorf("missing spec-url or file"))
	}

	var (
		reader io.Reader
		err    error
	)
	switch sch := g.URL.Scheme; sch {
	case "http", "https":
		var (
			cli = &http.Client{}
			rsp *http.Response
		)
		rsp, err = cli.Get(g.URL.String())
		if err == nil && rsp != nil {
			reader = rsp.Body
		}
	default:
		reader, err = os.Open(g.URL.Path)
	}

	if err != nil {
		panic(errors.Errorf("open spec failed: %v %s", err, g.URL))
	}
	if reader == nil {
		panic(errors.Errorf("open spec failed: %s", g.URL))
	}
	if err = json.NewDecoder(reader).Decode(g.Spec); err != nil {
		panic(errors.Wrap(err, "decode spec content"))
	}
}

func (g *Generator) Output(cwd string) {
	var (
		pkg  = stringsx.LowerSnakeCase(g.ServiceName)
		path = filepath.Join(cwd, pkg)
	)
	ctx := context.Background()
	if g.VendorImportByGoMod {
		ctx = WithVendorImports(ctx, VendorImportsByGoMod(cwd))
	}

	tg := NewTypeGen(
		g.ServiceName,
		codegen.NewFile(pkg, filepath.Join(path, "types.go")),
	)
	if err := tg.Gen(ctx, g.Spec); err != nil {
		panic("gen type.go: " + err.Error())
	}

	cg := NewClientGen(
		g.ServiceName, tg,
		codegen.NewFile(pkg, filepath.Join(path, "client.go")))
	if err := cg.Gen(ctx, g.Spec); err != nil {
		panic("gen client.go: " + err.Error())
	}

	og := NewOperationGen(
		g.ServiceName, tg,
		codegen.NewFile(pkg, filepath.Join(path, "operations.go")))
	if err := og.Gen(ctx, g.Spec); err != nil {
		panic("gen operation.go: " + err.Error())
	}

	log.Printf("generated client of %s into %s", g.ServiceName, color.MagentaString(path))
}
