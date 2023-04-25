package tls

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	"github.com/pkg/errors"
)

var gDefaultTlsConfig = &tls.Config{
	ClientAuth:         tls.NoClientCert,
	ClientCAs:          nil,
	InsecureSkipVerify: true,
}

type X509KeyPair struct {
	KeyPath string `json:""`
	CrtPath string `json:""`
	CaPath  string `json:""`
	Key     string `json:"key"`
	Crt     string `json:"crt"`
	Ca      string `json:"ca"`
	conf    *tls.Config
}

func (c *X509KeyPair) read() (key, crt, ca []byte, err error, empty bool) {
	if c.Key+c.Ca+c.Crt != "" {
		return []byte(c.Key), []byte(c.Crt), []byte(c.Ca), nil, false
	}
	if len(c.KeyPath+c.CrtPath+c.CaPath) == 0 {
		empty = true
		return
	}
	empty = false
	var content []byte
	content, err = os.ReadFile(c.KeyPath)
	if err != nil {
		return
	}
	key = content
	content, err = os.ReadFile(c.CrtPath)
	if err != nil {
		return
	}
	crt = content
	content, err = os.ReadFile(c.CaPath)
	if err != nil {
		return
	}
	ca = content
	return
}

func (c *X509KeyPair) Init() error {
	if c == nil {
		return nil
	}

	key, crt, ca, err, empty := c.read()
	if err != nil {
		return err
	}

	if empty {
		return nil
	}

	cert, err := tls.X509KeyPair(crt, key)
	if err != nil {
		return err
	}
	pool := x509.NewCertPool()
	ok := pool.AppendCertsFromPEM(ca)
	if !ok {
		return errors.Wrap(err, "failed to append cert")
	}
	c.conf = &tls.Config{
		RootCAs:            pool,
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}
	return nil
}

func (c *X509KeyPair) TLSConfig() *tls.Config {
	if c == nil {
		return gDefaultTlsConfig
	}
	return c.conf
}
