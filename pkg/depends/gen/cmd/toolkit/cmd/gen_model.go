package cmd

import (
	"github.com/spf13/cobra"

	"github.com/machinefi/w3bstream/pkg/depends/kit/modelgen"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
)

func init() {
	cmd := &cobra.Command{
		Use:   "model",
		Short: "generate database model",
		Run: func(cmd *cobra.Command, args []string) {
			if database == "" {
				panic("database required")
			}
			for _, arg := range args {
				run("model", func(pkg *pkgx.Pkg) Generator {
					g := modelgen.New(pkg)
					g.WithComments = true
					g.WithTableInterfaces = true
					g.WithMethods = true
					g.StructName = arg
					g.Database = database
					g.TableName = tableName

					g.Scan()
					return g
				}, arg)
			}
		},
	}
	cmd.Flags().StringVarP(
		&database, "database", "", "", "(required) database name",
	)
	cmd.Flags().StringVarP(
		&tableName, "table-name", "t", "", "custom table name",
	)
	cmd.Flags().BoolVarP(
		&withMethods, "with-methods", "", true, "with table methods",
	)

	Gen.AddCommand(cmd)
}

var (
	database    string
	tableName   string
	withMethods bool
)
