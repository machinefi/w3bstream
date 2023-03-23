package mqtt

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
)

type Client struct {
	cid     string // cid client id
	topic   string // topic registered topic
	qos     QOS    // qos should be 0, 1 or 2
	retain  bool
	timeout time.Duration //
	cli     mqtt.Client
}

func (c *Client) Cid() string { return c.cid }

func (c *Client) WithTopic(topic string) *Client {
	if c.topic == topic {
		return c
	}
	c2 := *c
	c2.topic = topic
	return &c2
}

func (c *Client) WithQoS(qos QOS) *Client {
	if c.qos == qos {
		return c
	}
	c2 := *c
	c2.qos = qos
	return &c2
}

func (c *Client) WithTimeout(timeout time.Duration) *Client {
	if c.timeout == timeout {
		return c
	}
	c2 := *c
	c2.timeout = timeout
	return &c2
}

func (c *Client) WithRetain(retain bool) *Client {
	if retain == c.retain {
		return c
	}
	c2 := *c
	c2.retain = retain
	return &c2
}

func (c *Client) connect() error {
	return c.wait(c.cli.Connect(), "connect")
}

func (c *Client) wait(tok mqtt.Token, act string) error {
	waited := false
	if c.timeout == 0 {
		waited = tok.Wait()
	} else {
		waited = tok.WaitTimeout(c.timeout)
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
	)
}

func (c *Client) Subscribe(cb mqtt.MessageHandler) error {
	if c.topic == "" {
		return errors.New("topic is empty")
	}
	return c.wait(
		c.cli.Subscribe(c.topic, byte(c.qos), cb),
		"sub",
	)
}

func (c *Client) Unsubscribe() error {
	if c.topic == "" {
		return errors.New("topic is empty")
	}
	return c.wait(
		c.cli.Unsubscribe(c.topic),
		"Unsub",
	)
}
