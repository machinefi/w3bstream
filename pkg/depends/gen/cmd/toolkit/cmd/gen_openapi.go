package cmd

import (
	"fmt"

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
				fmt.Println(validators)
				g := openapi.NewGenerator(pkg, "projectName", "^[a-z0-9_]{6,32}$")
				g.Scan(nil)
				return g
			})
		},
	}

	cmd.Flags().StringSliceVarP(&validators, "validators", "", nil, "import user defined validators")

	Gen.AddCommand(cmd)
}

var validators []string
