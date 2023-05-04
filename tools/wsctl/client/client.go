package client

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/tools/wsctl/config"
)

// Client defines the interface of an wsctl client
type Client interface {
	// Config returns the config of the client
	Config() *config.Config
	ConfigInfo() config.Info
	Address() string
	// SelectTranslation select a translation based on UILanguage
	SelectTranslation(map[config.Language]string) string
	// Call http call
	Call(req *http.Request) ([]byte, error)
	CallWithToken(req *http.Request, token string) ([]byte, error)
}

type client struct {
	cfg config.Info
}

// NewClient creates a new wsctl client
func NewClient(cfg config.Info) Client {
	return &client{
		cfg: cfg,
	}
}

func (c *client) Config() *config.Config {
	return c.cfg.Config()
}

func (c *client) ConfigInfo() config.Info {
	return c.cfg
}

func (c *client) Address() string {
	pvk, err := c.loadPrivateKey(c.Config().PrivateKey)
	if err != nil {
		panic(err)
	}
	return strings.ToLower(crypto.PubkeyToAddress(pvk.PublicKey).String())
}

func (c *client) loadPrivateKey(pvkStr string) (*ecdsa.PrivateKey, error) {
	var pvk *ecdsa.PrivateKey
	if len(pvkStr) > 0 {
		log.Println("loaded private key from the config.")
		pvk = crypto.ToECDSAUnsafe(common.FromHex(pvkStr))
	} else {
		var err error
		if pvk, err = crypto.GenerateKey(); err != nil {
			return nil, err
		}
		c.Config().PrivateKey = hex.EncodeToString(crypto.FromECDSA(pvk))
		if err := c.ConfigInfo().WriteConfig(); err != nil {
			return nil, err
		}
		log.Printf("no private key is found in the config; a new one %s is randomly generated.\n", c.Config().PrivateKey)
	}
	return pvk, nil
}

func (c *client) SelectTranslation(trls map[config.Language]string) string {
	trl, ok := trls[c.Config().Language]
	if !ok {
		log.Panic(errors.New("failed to pick a translation"))
	}
	return trl
}

func (c *client) Call(req *http.Request) ([]byte, error) {
	token, err := c.getToken()
	if err != nil {
		return nil, err
	}
	var resp []byte
	if resp, err = c.call(req, token); err != nil {
		c.Config().LoginToken = ""
		// retry with the new token by relogining
		token, err = c.getToken()
		if err != nil {
			return nil, err
		}
		if resp, err = c.call(req, token); err != nil {
			return nil, errors.Wrap(err, "failed to call w3bstream api")
		}
	}

	// dump the new token into config
	c.Config().LoginToken = token
	if err := c.ConfigInfo().WriteConfig(); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *client) CallWithToken(req *http.Request, token string) ([]byte, error) {
	return c.call(req, token)
}

func (c *client) call(req *http.Request, token string) ([]byte, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call w3bstream api")
	}
	defer resp.Body.Close()
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return nil, err
	}
	fmt.Println("\033[36mHttp response:\033[0m")
	fmt.Println(string(dump))
	if resp.StatusCode >= 400 {
		return nil, errors.New("error in the http response")
	}
	return io.ReadAll(resp.Body)
}

func (c *client) getToken() (string, error) {
	if len(c.Config().LoginToken) > 0 {
		log.Println("load login token in the config:", c.Config().LoginToken)
		return c.Config().LoginToken, nil
	}
	newtoken, err := c.login()
	if err != nil {
		return "", err
	}
	log.Println("new token!", newtoken)
	return newtoken, nil
}
