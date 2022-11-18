package client

import (
	"errors"

	"github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/wsctl/config"
)

// Client defines the interface of an wsctl client
type Client interface {
	// Config returns the config of the client
	Config() config.Config
	// ConfigFilePath returns the file path of the config
	ConfigFilePath() string
	// SelectTranslation select a translation based on UILanguage
	SelectTranslation(map[config.Language]string) string
}

type client struct {
	cfg            config.Config
	configFilePath string
	logger         log.Logger
}

// NewClient creates a new wsctl client
func NewClient(cfg config.Config, configFilePath string, logger log.Logger) Client {
	return &client{
		cfg:            cfg,
		configFilePath: configFilePath,
		logger:         logger,
	}
}

func (c *client) Config() config.Config {
	return c.cfg
}

// ConfigFilePath returns the file path for the config.
func (c *client) ConfigFilePath() string {
	return c.configFilePath
}

func (c *client) SelectTranslation(trls map[config.Language]string) string {
	trl, ok := trls[c.cfg.Language]
	if !ok {
		c.logger.Panic(errors.New("failed to pick a translation"))
	}
	return trl
}
