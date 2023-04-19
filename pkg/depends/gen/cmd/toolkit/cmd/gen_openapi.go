package cmd

import (
	"github.com/spf13/cobra"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransportgen/openapi"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
)

func init() {
	cmd := &cobra.Command{
		Use:   "openapi",
		Short: "generate openapi spec for current project",
		Run: func(cmd *cobra.Command, args []string) {
			run("openapi", func(pkg *pkgx.Pkg) Generator {
				g := openapi.NewGenerator(pkg)
				g.Scan(nil)
				return g
			})
		},
	}

	Gen.AddCommand(cmd)
}
