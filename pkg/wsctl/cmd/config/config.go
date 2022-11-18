package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/machinefi/w3bstream/pkg/wsctl/client"
	"github.com/machinefi/w3bstream/pkg/wsctl/config"
)

const (
	_ipPattern       = `((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)`
	_domainPattern   = `[a-zA-Z0-9][a-zA-Z0-9_-]{0,62}(\.[a-zA-Z0-9][a-zA-Z0-9_-]{0,62})*(\.[a-zA-Z][a-zA-Z0-9]{0,10}){1}`
	_localPattern    = "localhost"
	_endpointPattern = "(" + _ipPattern + "|(" + _domainPattern + ")" + "|(" + _localPattern + "))" + `(:\d{1,5})?`

	_defaultConfigFileName = "config.default"
	_defaultEndpoint       = "http://localhost:8888"
)

var (
	_configDir       = os.Getenv("HOME") + "/.config/wsctl/default"
	_endpointCompile = regexp.MustCompile("^" + _endpointPattern + "$")
)

// Multi-language support
var (
	_configCmdShorts = map[config.Language]string{
		config.English: "Manage the configuration of wsctl",
		config.Chinese: "wsctl 配置管理",
	}
)

// NewConfigCmd represents the new node command.
func NewConfigCmd(client client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: client.SelectTranslation(_configCmdShorts),
	}
	cmd.AddCommand(newConfigSetCmd(client))
	cmd.AddCommand(newConfigGetCmd(client))
	cmd.AddCommand(newConfigResetCmd(client))

	return cmd
}

// InitConfig load config data from default config file
func InitConfig() (config.Config, string, error) {
	info := &info{
		config: config.Config{},
	}

	// Create path to config directory
	err := os.MkdirAll(_configDir, 0700)
	if err != nil {
		return info.config, info.configFile, err
	}
	info.configFile = filepath.Join(_configDir, _defaultConfigFileName)

	// Load or reset config file
	err = info.loadConfig()
	if os.IsNotExist(err) {
		err = info.reset()
	}
	if err != nil {
		return info.config, info.configFile, err
	}

	// Check completeness of config file
	completeness := true
	if info.config.Language == "" {
		info.config.Language = config.English
		completeness = false
	}
	if info.config.Endpoint == "" {
		info.config.Endpoint = _defaultEndpoint
		completeness = false
	}
	if !completeness {
		if err = info.writeConfig(); err != nil {
			return info.config, info.configFile, err
		}
	}
	if !isSupportedLanguage(info.config.Language) {
		fmt.Printf("Warn: Language %s is not supported, using English.\n", info.config.Language)
	}
	return info.config, info.configFile, nil
}

// info contains the information of config file
type info struct {
	config     config.Config
	configFile string // Path to config file
}

// loadConfig loads config file in yaml format
func (c *info) loadConfig() error {
	in, err := os.ReadFile(c.configFile)
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal(in, &c.config); err != nil {
		return errors.Wrap(err, "failed to unmarshal config")
	}
	return nil
}

// reset resets all values of config
func (c *info) reset() error {
	c.config.Endpoint = _defaultEndpoint
	c.config.Language = config.English

	err := c.writeConfig()
	if err != nil {
		return err
	}

	fmt.Println("Config set to default values")
	return nil
}

// writeConfig writes to config file
func (c *info) writeConfig() error {
	out, err := yaml.Marshal(&c.config)
	if err != nil {
		return errors.Wrap(err, "failed to marshal config")
	}
	if err := os.WriteFile(c.configFile, out, 0600); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to write to config file %s", c.configFile))
	}
	return nil
}

func isSupportedLanguage(l config.Language) bool {
	for _, lang := range config.SupportedLanguage {
		if strings.EqualFold(string(l), string(lang)) {
			return true
		}
	}
	return false
}

// isValidEndpoint makes sure the endpoint matches the endpoint match pattern
func isValidEndpoint(endpoint string) bool {
	return _endpointCompile.MatchString(endpoint)
}

// jsonString returns json string for message
func jsonString(input interface{}) (string, error) {
	byteAsJSON, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		return "", errors.Wrap(err, "failed to JSON marshal config field")
	}
	return fmt.Sprint(string(byteAsJSON)), nil
}
