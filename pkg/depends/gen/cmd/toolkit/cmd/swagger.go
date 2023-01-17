package cmd

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/httpswaggergen"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
	"github.com/spf13/cobra"
)

var Swagger = &cobra.Command{
	Use:   "swagger",
	Short: "scan current project and generate openapi.json",
	Run: func(cmd *cobra.Command, args []string) {
		run("openapi", func(pkg *pkgx.Pkg) Generator {
			g := httpswaggergen.NewOpenAPIGenerator(pkg)
			g.Scan(nil)
			return g
		})
	},
}
