package cmd

import (
	"github.com/spf13/cobra"

	"github.com/machinefi/w3bstream/pkg/depends/kit/enumgen"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
)

func init() {
	cmd := &cobra.Command{
		Use:   "enum",
		Short: "generate interfaces of enumeration",
		Run: func(cmd *cobra.Command, args []string) {
			run("enum", func(pkg *pkgx.Pkg) Generator {
				g := enumgen.New(pkg)
				g.Scan(args...)
				return g
			}, args...)
		},
	}

	Gen.AddCommand(cmd)
}
