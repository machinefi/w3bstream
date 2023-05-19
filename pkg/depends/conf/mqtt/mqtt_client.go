package mqtt

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
)

type Client struct {
	cid         string // cid client id
	topic       string // topic registered topic
	qos         QOS    // qos should be 0, 1 or 2
	retain      bool
	pubTimeout  time.Duration
	subTimeout  time.Duration
	connTimeout time.Duration
	cli         mqtt.Client
}

func (c *Client) Cid() string { return c.cid }

func (c *Client) WithTopic(topic string) *Client {
	c2 := *c
	c2.topic = topic
	return &c2
}

func (c *Client) WithQoS(qos QOS) *Client {
	c2 := *c
	c2.qos = qos
	return &c2
}

func (c *Client) WithSubTimeout(timeout time.Duration) *Client {
	c2 := *c
	c2.subTimeout = timeout
	return &c2
}

func (c *Client) WithConnTimeout(timeout time.Duration) *Client {
	c2 := *c
	c2.connTimeout = timeout
	return &c2
}

func (c *Client) WithRetain(retain bool) *Client {
	c2 := *c
	c2.retain = retain
	return &c2
}

func (c *Client) connect() error {
	return c.wait(c.cli.Connect(), "connect", c.connTimeout)
}

func (c *Client) wait(tok mqtt.Token, act string, timeout time.Duration) error {
	waited := false
	if timeout == 0 {
		waited = tok.Wait()
	} else {
		waited = tok.WaitTimeout(timeout)
	}
	if !waited {
		return errors.New(act + " timeout")
	}
	if err := tok.Error(); err != nil {
		return errors.Wrap(err, act+" error")
	}
	return nil
}

func (c *Client) Publish(payload interface{}) error {
	if c.topic == "" {
		return errors.New("topic is empty")
	}
	return c.wait(
		c.cli.Publish(c.topic, byte(c.qos), c.retain, payload),
		"pub",
		c.pubTimeout,
	)
}

func (c *Client) Subscribe(handler mqtt.MessageHandler) error {
	if c.topic == "" {
		return errors.New("topic is empty")
	}
	return c.wait(
		c.cli.Subscribe(c.topic, byte(c.qos), handler),
		"sub",
		c.subTimeout,
	)
}

func (c *Client) Unsubscribe() error {
	if c.topic == "" {
		return errors.New("topic is empty")
	}
	return c.wait(
		c.cli.Unsubscribe(c.topic),
		"Unsub",
		c.subTimeout,
	)
}
