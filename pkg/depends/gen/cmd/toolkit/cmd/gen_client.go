package cmd

import (
	neturl "net/url"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransportgen/client"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
	"github.com/spf13/cobra"
)

var (
	oasURL string
	name   string
)

func init() {
	Gen.AddCommand(cmd)

	cmd.Flags().StringVarP(&oasURL, "url", "", "", "client spec url")
	cmd.Flags().StringVarP(&name, "name", "", "", "service name")
}

var cmd = &cobra.Command{
	Use:     "client",
	Example: "client demo",
	Short:   "generate client by open api",
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
