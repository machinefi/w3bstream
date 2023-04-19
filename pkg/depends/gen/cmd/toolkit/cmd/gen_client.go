package cmd

import (
	neturl "net/url"

	"github.com/spf13/cobra"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransportgen/client"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
)

var (
	oasURL string
	name   string
)

func init() {
	cmd := &cobra.Command{
		Use:     "client",
		Example: "client demo",
		Short:   "generate client by openapi spec",
		Run: func(cmd *cobra.Command, args []string) {
			if name == "" {
				cmd.Println("require service name")
				return
			}
			u, err := neturl.Parse(oasURL)
			if err != nil {
				cmd.Printf("invalid url: %s", oasURL)
				return
			}

			run("client", func(pkg *pkgx.Pkg) Generator {
				g := client.NewGenerator(name, u)
				g.Load()
				return g
			})
		},
	}

	cmd.Flags().StringVarP(&oasURL, "url", "", "", "client spec url")
	cmd.Flags().StringVarP(&name, "name", "", "", "service name")

	Gen.AddCommand(cmd)
}
