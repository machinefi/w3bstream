package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/pkg/errors"
	"golang.org/x/term"

	"github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/wsctl/config"
)

// Multi-language support
var (
	_userNameInput = map[config.Language]string{
		config.English: "Please enter your name: ",
		config.Chinese: "请输入用户名：",
	}
	_userPasswordInput = map[config.Language]string{
		config.English: "Please enter password: ",
		config.Chinese: "请输入密码：",
	}
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
	Call(url string, req *http.Request) (*http.Response, error)
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

type callResp struct {
	Code uint64 `json:"code"`
}

func (c *client) Call(url string, req *http.Request) (*http.Response, error) {
	resp, err := c.call(url, req, c.getToken())
	if err != nil {
		return nil, errors.Wrap(err, "failed to call w3bstream api")
	}
	cr := callResp{}
	if err := json.NewDecoder(resp.Body).Decode(&cr); err != nil {
		errors.Wrap(err, "failed to decode w3bstream api responce")
	}
	if cr.Code/1e6 == 401 {
		c.login()
		return c.call(url, req, c.getToken())
	}
	return resp, err
}

func (c *client) call(url string, req *http.Request, token string) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	cli := &http.Client{}
	return cli.Do(req)
}

func (c *client) getToken() string {
	if t := c.token.Load(); t != nil {
		return t.(string)
	}
	c.login()
	return c.token.Load().(string)
}

type tokenResp struct {
	Token string `json:"token"`
}

func (c *client) login() {
	var userName, password string
	fmt.Println(c.SelectTranslation(_userNameInput))
	if _, err := fmt.Scanln(&userName); err != nil {
		c.logger.Panic(errors.Wrap(err, "failed to read username"))
	}

	fmt.Println(c.SelectTranslation(_userPasswordInput))
	p, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		c.logger.Panic(errors.Wrap(err, "failed to read password"))
	}
	password = string(p)

	body := fmt.Sprintf(`{"username":"%s","password":"%s"}`, userName, password)
	url := fmt.Sprintf("%s/srv-applet-mgr/v0/login", c.cfg.Endpoint)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		c.logger.Panic(errors.Wrap(err, "failed to create login request"))
	}
	req.Header.Set("Content-Type", "application/json")
	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		c.logger.Panic(errors.Wrapf(err, "failed to login %s", url))
	}
	defer resp.Body.Close()

	ts := tokenResp{}

	if err := json.NewDecoder(resp.Body).Decode(&ts); err != nil {
		c.logger.Panic(errors.Wrap(err, "failed to decode login response"))
	}
	c.token.Store(ts.Token)
}
