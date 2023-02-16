package mqtt

import (
	"crypto/tls"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/x/mapx"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/retry"
)

type Broker struct {
	Server        types.Endpoint `json:"broker,string"`
	Retry         retry.Retry    `json:"-"`
	Timeout       types.Duration `json:"-"`
	Keepalive     types.Duration `json:"-"`
	RetainPublish bool           `json:"retain"`
	QoS           QOS            `json:"-"`

	agents *mapx.Map[string, *Client]
}

func (b *Broker) SetDefault() {
	b.Retry.SetDefault()
	if b.Keepalive == 0 {
		b.Keepalive = types.Duration(3 * time.Hour)
	}
	if b.Server.IsZero() {
		b.Server.Hostname, b.Server.Port = "127.0.0.1", 1883
	}
	b.Server.Scheme = "mqtt"
	if b.agents == nil {
		b.agents = mapx.New[string, *Client]()
	}
}

func (b *Broker) Init() error {
	return b.Retry.Do(func() error {
		cid := uuid.New().String()
		_, err := b.Client(cid)
		if err != nil {
			return err
		}
		b.Close(cid)
		return nil
	})
}

func (b *Broker) options() *mqtt.ClientOptions {
	opt := mqtt.NewClientOptions()
	if !b.Server.IsZero() {
		opt = opt.AddBroker(b.Server.SchemeHost())
	}
	if b.Server.Username != "" {
		opt.SetUsername(b.Server.Username)
		if b.Server.Password != "" {
			opt.SetPassword(b.Server.Password.String())
		}
	}

	opt.SetKeepAlive(b.Keepalive.Duration())
	opt.SetWriteTimeout(b.Timeout.Duration())
	opt.SetConnectTimeout(b.Timeout.Duration())
	return opt
}

func (b *Broker) Client(cid string) (*Client, error) {
	opt := b.options()
	if cid != "" {
		opt.SetClientID(cid)
	}
	// TODO support TLS
	if b.Server.IsTLS() {
		opt.SetTLSConfig(&tls.Config{
			ClientAuth:         tls.NoClientCert,
			ClientCAs:          nil,
			InsecureSkipVerify: true,
		})
	}
	if b.Server.Username != "" {
		opt.SetUsername(b.Server.Username)
		opt.SetPassword(b.Server.Password.String())
	}
	return b.ClientWithOptions(cid, opt)
}

func (b *Broker) ClientWithOptions(cid string, opt *mqtt.ClientOptions) (*Client, error) {
	client, err := b.agents.LoadOrStore(
		cid,
		func() (*Client, error) {
			if opt.WriteTimeout == 0 {
				opt.WriteTimeout = 10 * time.Second
			}
			if opt.ConnectTimeout == 0 {
				opt.ConnectTimeout = 10 * time.Second
			}
			c := &Client{
				cid:    cid,
				qos:    b.QoS,
				retain: b.RetainPublish,
				cli:    mqtt.NewClient(opt),
			}
			if err := c.connect(); err != nil {
				return nil, err
			}
			return c, nil
		},
	)
	if err != nil {
		return nil, err
	}
	if !client.cli.IsConnectionOpen() && !client.cli.IsConnected() {
		b.agents.Remove(cid)
		return b.Client(cid)
	}
	return client, nil
}

func (b *Broker) Close(cid string) {
	if c, ok := b.agents.LoadAndRemove(cid); ok && c != nil {
		c.cli.Disconnect(500)
	}
}

func init() {
	// mqtt.ERROR = log.New(os.Stderr, "******************mqtt**[ERR]", 0)
	// mqtt.CRITICAL = log.New(os.Stderr, "******************mqtt**[CRI]", 0)
	// mqtt.WARN = log.New(os.Stderr, "******************mqtt**[WAR]", 0)
	// mqtt.DEBUG = log.New(os.Stderr, "******************mqtt**", 0)
}
