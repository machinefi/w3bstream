package client

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"sync/atomic"

	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/tools/wsctl/config"
)

// Client defines the interface of an wsctl client
type Client interface {
	// Config returns the config of the client
	Config() config.Config
	// ConfigFilePath returns the file path of the config
	ConfigFilePath() string
	// SelectTranslation select a translation based on UILanguage
	SelectTranslation(map[config.Language]string) string
	// Call http call
	Call(req *http.Request) ([]byte, error)
}

type client struct {
	cfg            config.Config
	configFilePath string
	logger         log.Logger
	token          atomic.Value
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

func (c *client) Call(req *http.Request) ([]byte, error) {
	token, err := c.getToken()
	if err != nil {
		return nil, err
	}
	resp, err := c.call(req, token)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call w3bstream api")
	}
	defer resp.Body.Close()
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return nil, err
	}
	fmt.Println("\033[36mhttp response:\033[0m")
	fmt.Println(string(dump))
	if resp.StatusCode >= 400 {
		return nil, errors.New("error in the http response")
	}
	return io.ReadAll(resp.Body)
}

func (c *client) call(req *http.Request, token string) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	cli := &http.Client{}
	return cli.Do(req)
}

func (c *client) getToken() (string, error) {
	if t := c.token.Load(); t != nil {
		return t.(string), nil
	}
	if err := c.login(); err != nil {
		return "", err
	}
	return c.token.Load().(string), nil
}
