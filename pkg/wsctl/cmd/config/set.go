package config

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/machinefi/w3bstream/pkg/wsctl/client"
	"github.com/machinefi/w3bstream/pkg/wsctl/config"
)

var _validArgs = []string{"endpoint", "language"}

var (
	_configSetUse = map[config.Language]string{
		config.English: "set VARIABLE VALUE",
		config.Chinese: "set 变量 值",
	}
	_configSetUseCmdShorts = map[config.Language]string{
		config.English: "Set config fields for wsctl",
		config.Chinese: "为 wsctl 设置配置字段",
	}
	_configSetUseCmdLong = map[config.Language]string{
		config.English: "Set config fields for wsctl\nValid Variables: [" + strings.Join(_validArgs, ", ") + "]",
		config.Chinese: "为 wsctl 设置配置字段\n有效变量: [" + strings.Join(_validArgs, ", ") + "]",
	}
)

// newConfigSetCmd is a command to set config fields from wsctl
func newConfigSetCmd(client client.Client) *cobra.Command {
	return &cobra.Command{
		Use:       client.SelectTranslation(_configSetUse),
		Short:     client.SelectTranslation(_configSetUseCmdShorts),
		Long:      client.SelectTranslation(_configSetUseCmdLong),
		ValidArgs: _validArgs,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("accepts 2 arg(s), received %d\n"+
					"Valid arg(s): %s", len(args), _validArgs)
			}
			return cobra.OnlyValidArgs(cmd, args[:1])
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			info := info{config: client.Config(), configFile: client.ConfigFilePath()}
			result, err := info.set(args)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("problem setting config fields %+v", args))
			}
			cmd.Println(result)
			return nil
		},
	}
}

// set sets config variable
func (c *info) set(args []string) (string, error) {
	switch args[0] {
	case "endpoint":
		if !isValidEndpoint(args[1]) {
			return "", errors.Errorf("endpoint %s is not valid", args[1])
		}
		c.config.Endpoint = args[1]
	case "language":
		if !isSupportedLanguage(config.Language(args[1])) {
			return "", errors.Errorf("language %s is not supported\nSupported languages: %s",
				args[1], config.SupportedLanguage)
		}
		c.config.Language = config.Language(args[1])
	default:
		return "", ErrConfigNotMatch
	}

	if err := c.writeConfig(); err != nil {
		return "", err
	}

	return cases.Title(language.Und).String(args[0]) + " is set to " + args[1], nil
}
