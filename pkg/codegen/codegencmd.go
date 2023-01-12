package codegen

import (
	"errors"

	"github.com/machinefi/w3bstream/pkg/codegen/operators"
	"github.com/spf13/cobra"
)

func NewCodeGenCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "codegen configfile",
		Short: "generate code for given config file",
		Long:  "generate code based on the given config file in yaml",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("invalid parameter")
			}
			return cobra.OnlyValidArgs(cmd, args[:0])
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			cmd.Println(args)
			builder := NewObservableBuilder()
			builder.Append(operators.NewFilterGenerator("wasm_func_name"))
			code, err := builder.Build()
			if err != nil {
				return err
			}
			cmd.Println(code)
			// TODO: parse args[0] as config file in yaml
			return nil
		},
	}
}
