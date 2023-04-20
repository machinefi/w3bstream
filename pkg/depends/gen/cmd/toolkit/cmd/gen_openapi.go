package cmd

import (
	"encoding/json"

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
				nameRulePairs := make([]string, 0)
				if err := json.Unmarshal([]byte(validators), &nameRulePairs); err != nil {
					panic("should pass validators by name rule pairs with a json array format")
				}
				g := openapi.NewGenerator(pkg, nameRulePairs...)
				g.Scan(nil)
				return g
			})
		},
	}

	cmd.Flags().StringVarP(&validators, "validators", "", "", "import user defined validators")

	Gen.AddCommand(cmd)
}

var validators string
