package config

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/machinefi/w3bstream/tools/wsctl/client"
	"github.com/machinefi/w3bstream/tools/wsctl/config"
)

var _validGetArgs = []string{"endpoint", "language", "all"}

var (
	_configGetUse = map[config.Language]string{
		config.English: "get VARIABLE",
		config.Chinese: "get 变量",
	}
	_configGetUseCmdShorts = map[config.Language]string{
		config.English: "Get config fields from wsctl",
		config.Chinese: "从 wsctl 获取配置字段",
	}
	_configGetUseCmdLong = map[config.Language]string{
		config.English: "Get config fields from wsctl\nValid Variables: [" + strings.Join(_validGetArgs, ", ") + "]",
		config.Chinese: "从 wsctl 获取配置字段\n有效变量: [" + strings.Join(_validGetArgs, ", ") + "]",
	}
)

// newConfigGetCmd is a command to get config fields from wstcl.
func newConfigGetCmd(client client.Client) *cobra.Command {
	return &cobra.Command{
		Use:       client.SelectTranslation(_configGetUse),
		Short:     client.SelectTranslation(_configGetUseCmdShorts),
		Long:      client.SelectTranslation(_configGetUseCmdLong),
		ValidArgs: _validGetArgs,
		Args:      cobra.ExactValidArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			result, err := get(client.Config(), args[0])
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("issue fetching config value %s", args[0]))
			}
			cmd.Println(result)
			return nil
		},
	}
}

// get retrieves a config item from its key.
func get(c *config.Config, arg string) (string, error) {
	switch arg {
	case "endpoint":
		if c.Endpoint == "" {
			return "", ErrEmptyEndpoint
		}
		return fmt.Sprintf("%s", c.Endpoint), nil
	case "language":
		return string(c.Language), nil
	case "all":
		return jsonString(c)
	default:
		return "", ErrConfigNotMatch
	}
}
