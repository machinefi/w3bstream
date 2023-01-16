package cmd

import (
	"github.com/spf13/cobra"

	"github.com/machinefi/w3bstream/pkg/depends/kit/validator/strfmtgen"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
)

func init() {
	cmd := &cobra.Command{
		Use:   "strfmt",
		Short: "generate strfmt validator",
		Run: func(cmd *cobra.Command, args []string) {
			if filename == "" {
				panic("filename required")
			}
			run("strfmt", func(pkg *pkgx.Pkg) Generator {
				return strfmtgen.NewGenerator(pkg, filename)
			})
		},
	}
	cmd.Flags().StringVarP(&filename, "filename", "f", "", "(required) filename")

	Gen.AddCommand(cmd)
}

var filename string
