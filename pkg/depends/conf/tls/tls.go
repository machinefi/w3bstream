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
	KeyPath string `json:"-"`
	CrtPath string `json:"-"`
	CaPath  string `json:"-"`
	Key     []byte `env:"-" json:"key"`
	Crt     []byte `env:"-" json:"crt"`
	Ca      []byte `env:"-" json:"ca"`
	conf    *tls.Config
}

func (c *X509KeyPair) Init() error {
	if c == nil {
		return nil
	}
	if c.KeyPath != "" {
		content, err := os.ReadFile(c.KeyPath)
		if err != nil {
			return err
		}
		c.Key = content
	}
	if c.CrtPath != "" {
		content, err := os.ReadFile(c.CrtPath)
		if err != nil {
			return err
		}
		c.Crt = content
	}
	if c.CaPath != "" {
		content, err := os.ReadFile(c.CaPath)
		if err != nil {
			return err
		}
		c.Ca = content
	}

	if len(c.Crt) == 0 || len(c.Key) == 0 || len(c.Ca) == 0 {
		c.conf = gDefaultTlsConfig
		return nil
	}
	cert, err := tls.X509KeyPair(c.Crt, c.Key)
	if err != nil {
		return err
	}
	pool := x509.NewCertPool()
	ok := pool.AppendCertsFromPEM(c.Ca)
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
