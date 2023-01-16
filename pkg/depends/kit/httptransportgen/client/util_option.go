package client

import (
	"context"

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
