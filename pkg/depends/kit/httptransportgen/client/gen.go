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
	"github.com/pkg/errors"
	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/packages"

	"github.com/machinefi/w3bstream/pkg/depends/gen/codegen"
	"github.com/machinefi/w3bstream/pkg/depends/oas"
	"github.com/machinefi/w3bstream/pkg/depends/x/stringsx"
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

	ctx := WithVendorImports(context.Background(), g.VendorImportsByGoMod(cwd))

	{
		f := codegen.NewFile(pkg, filepath.Join(path, "client.go"))
		err := NewClientGen(g.ServiceName, f).Gen(ctx, g.Spec)
		if err != nil {
			panic("gen client.go: " + err.Error())
		}
	}

	{
		f := codegen.NewFile(pkg, filepath.Join(path, "operations.go"))
		err := NewOperationGen(g.ServiceName, f).Gen(ctx, g.Spec)
		if err != nil {
			panic("gen operation.go: " + err.Error())
		}
	}

	{
		f := codegen.NewFile(pkg, filepath.Join(path, "types.go"))
		err := NewTypeGen(g.ServiceName, f).Gen(ctx, g.Spec)
		if err != nil {
			panic("gen type.go: " + err.Error())
		}
	}

	log.Printf("generated client of %s into %s", g.ServiceName, color.MagentaString(path))
}

func (g *Generator) VendorImportsByGoMod(cwd string) map[string]bool {
	imports := map[string]bool{}

	if g.VendorImportByGoMod {
		pkgs, err := packages.Load(nil, "std")
		if err != nil {
			panic(err)
		}

		for _, p := range pkgs {
			imports[p.PkgPath] = true
		}

		for d := cwd; d != "/"; d = filepath.Join(d, "../") {
			fgomod := filepath.Join(d, "go.mod")

			if data, err := os.ReadFile(fgomod); err != nil {
				if !os.IsNotExist(err) {
					panic(err)
				}
			} else {
				f, _ := modfile.Parse(fgomod, data, nil)

				imports[f.Module.Mod.Path] = true

				for _, r := range f.Require {
					imports[r.Mod.Path] = true
				}

				break
			}

			d = filepath.Join(d, "../")
		}
	}

	return imports
}
