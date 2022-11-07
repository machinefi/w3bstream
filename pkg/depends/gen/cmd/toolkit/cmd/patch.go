package cmd

import "github.com/spf13/cobra"

var Patch = &cobra.Command{
	Use:   "patch",
	Short: "patch code to go root",
}

func init() {
	cmd := &cobra.Command{
		Use:   "goid",
		Short: "patch runtime goid for debug",
		Run: func(cmd *cobra.Command, args []string) {
			patchGoID()
		},
	}

	Patch.AddCommand(cmd)
}
