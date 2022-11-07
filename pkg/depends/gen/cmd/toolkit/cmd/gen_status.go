package cmd

import (
	"github.com/spf13/cobra"

	"github.com/machinefi/w3bstream/pkg/depends/kit/statusxgen"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
)

func init() {
	cmd := &cobra.Command{
		Use:     "status",
		Aliases: []string{"status-error", "error"},
		Short:   "generate interfaces of status error",
		Run: func(cmd *cobra.Command, args []string) {
			run("status", func(pkg *pkgx.Pkg) Generator {
				g := statusxgen.New(pkg)
				g.Scan(args...)
				return g
			}, args...)
		},
	}

	Gen.AddCommand(cmd)
}
