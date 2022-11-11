package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/machinefi/w3bstream/pkg/wsctl/client"
	"github.com/machinefi/w3bstream/pkg/wsctl/config"
)

var (
	_configResetCmdShorts = map[config.Language]string{
		config.English: "Reset config to default",
		config.Chinese: "将配置重置为默认值",
	}
)

// newConfigResetCmd resets the config to the default values
func newConfigResetCmd(client client.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "reset",
		Short: client.SelectTranslation(_configResetCmdShorts),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			info := info{config: client.Config(), configFile: client.ConfigFilePath()}
			err := info.reset()
			if err != nil {
				return errors.Wrap(err, "failed to reset config")
			}
			cmd.Println("successfully reset config")
			return nil
		},
	}
}
