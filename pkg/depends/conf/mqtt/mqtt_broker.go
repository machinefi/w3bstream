package mqtt

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	conftls "github.com/machinefi/w3bstream/pkg/depends/conf/tls"
	"github.com/machinefi/w3bstream/pkg/depends/x/mapx"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/retry"
)

type Broker struct {
	Server        types.Endpoint
	Retry         retry.Retry
	Timeout       types.Duration
	Keepalive     types.Duration
	RetainPublish bool
	QoS           QOS
	Cert          *conftls.X509KeyPair

	agents *mapx.Map[string, *Client]
}

func (b *Broker) SetDefault() {
	b.Retry.SetDefault()
	if b.Keepalive == 0 {
		b.Keepalive = types.Duration(3 * time.Hour)
	}
	if b.Timeout == 0 {
		b.Timeout = types.Duration(10 * time.Second)
	}
	if b.Server.IsZero() {
		b.Server.Hostname, b.Server.Port = "127.0.0.1", 1883
	}
	if b.Server.Scheme == "" {
		b.Server.Scheme = "mqtt"
	}
	if b.agents == nil {
		b.agents = mapx.New[string, *Client]()
	}
	if b.QoS > QOS__ONLY_ONCE || b.QoS < 0 {
		b.QoS = QOS__ONCE
	}
}

func (b *Broker) Init() error {
	if b.Cert != nil {
		if err := b.Cert.Init(); err != nil {
			return err
		}
	}
	return b.Retry.Do(func() error {
		cid := uuid.NewString()
		defer b.CloseByCid(cid)
		_, err := b.Client(cid)
		if err != nil {
			return err
		}
		return nil
	})
}

func (b *Broker) options(cid string) *mqtt.ClientOptions {
	opt := mqtt.NewClientOptions()
	if cid == "" {
		cid = uuid.NewString()
	}
	opt.SetClientID(cid)
	if !b.Server.IsZero() {
		opt = opt.AddBroker(b.Server.String())
	}
	if b.Server.Username != "" {
		opt.SetUsername(b.Server.Username)
		if b.Server.Password != "" {
			opt.SetPassword(b.Server.Password.String())
		}
	}
	if b.Server.IsTLS() {
		opt.SetTLSConfig(b.Cert.TLSConfig())
	}

	opt.SetKeepAlive(b.Keepalive.Duration())
	opt.SetWriteTimeout(b.Timeout.Duration())
	opt.SetConnectTimeout(b.Timeout.Duration())
	return opt
}

func (b *Broker) Name() string { return "mqtt-broker-cli" }

func (b *Broker) LivenessCheck() map[string]string {
	m := map[string]string{}
	cid := uuid.NewString()
	defer b.CloseByCid(cid)
	if _, err := b.Client(cid); err != nil {
		m[b.Server.Host()] = err.Error()
		return m
	}
	m[b.Server.Host()] = "ok"
	return m
}

func (b *Broker) Client(cid string) (*Client, error) {
	return b.ClientWithOptions(b.options(cid))
}

func (b *Broker) ClientWithOptions(opt *mqtt.ClientOptions) (*Client, error) {
	return b.agents.LoadOrStore(
		opt.ClientID,
		func() (*Client, error) {
			c := &Client{
				cid:    opt.ClientID,
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
}

func (b *Broker) Close(c *Client) {
	b.CloseByCid(c.cid)
}

func (b *Broker) CloseByCid(cid string) {
	if c, ok := b.agents.LoadAndRemove(cid); ok && c != nil {
		c.cli.Disconnect(500)
	}
}
