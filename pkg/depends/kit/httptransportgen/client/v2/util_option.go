package client

import (
	"context"
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/packages"

	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
)

func OptionVendorImportByGoMod() WithOption {
	return func(o *Option) {
		o.VendorImportByGoMod = true
	}
}

type WithOption = func(o *Option)

type Option struct {
	VendorImportByGoMod bool `name:"vendor-import-by-go-mod" usage:"when enable vendor only import pkg exists in go mod"`
}

type ctxVendorImports struct{}

func WithVendorImports(ctx context.Context, imports map[string]bool) context.Context {
	return contextx.WithValue(ctx, ctxVendorImports{}, imports)
}

func VendorImportsFromContext(ctx context.Context) map[string]bool {
	if v, ok := ctx.Value(ctxVendorImports{}).(map[string]bool); ok {
		return v
	}
	return map[string]bool{}
}

func VendorImportsByGoMod(cwd string) map[string]bool {
	imports := map[string]bool{}

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
	return imports
}
