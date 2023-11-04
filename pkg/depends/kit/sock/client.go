package sock

import (
	"net"
	"time"

	qsche "github.com/machinefi/w3bstream/pkg/depends/kit/sche"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sock/sock_buf/packet"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sock/sock_buf/stream"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sock/sock_msg"
)

type Client struct {
	node *Node
	*ClientOption
}

func NewClient(options ...ClientOptionSetter) (*Client, error) {
	var (
		err error
		cli = &Client{
			ClientOption: &ClientOption{
				timeout:        5 * time.Second,
				readBufferCap:  DefaultBufferCapacity,
				writeBufferCap: DefaultBufferCapacity,
				workerPoolSize: 1024,
			},
		}
		n = NewNode()
	)

	for _, opt := range options {
		opt(cli.ClientOption)
	}

	if cli.remote == "" {
		return nil, ENodeInvalidRemoteAddr
	}
	n.c, err = net.DialTimeout(cli.protocol.Network(), cli.remote, cli.timeout)
	if err != nil {
		return nil, ENodeDial.WithError(err)
	}
	err = SetConnOption(n.c, cli.readBufferCap, cli.writeBufferCap)
	if err != nil {
		return nil, ENodeOption.WithError(err)
	}
	if cli.parser == nil {
		return nil, ENodeInvalidParser
	}
	n.SetParser(cli.parser)
	if cli.protocol == ProtocolUnknown {
		return nil, ENodeInvalidProtocol
	}
	n.protocol = cli.protocol
	if cli.protocol == ProtocolTCP {
		n.SetBuffers(stream.New(cli.readBufferCap), packet.New(cli.writeBufferCap))
	} else {
		n.SetBuffers(packet.New(cli.readBufferCap), packet.New(cli.writeBufferCap))
	}
	n.id = cli.nodeID
	if cli.routes != nil || cli.handler != nil {
		n.worker = qsche.RunConScheduler(cli.workerPoolSize)
	}
	n.routes = cli.routes
	n.handler = cli.handler

	n.SetParser(cli.parser)

	n.rq = make(chan sock_msg.Message, 1024)
	n.sq = make(chan sock_msg.Message, 1024)

	n.Run()
	cli.node = n

	if cli.onConnected != nil {
		cli.onConnected(n)
	}

	return cli, nil
}

func (c *Client) RecvMessage() (sock_msg.Message, error) {
	return c.node.ReadMessage(c.timeout)
}

func (c *Client) SendMessage(msg sock_msg.Message) (err error) {
	return c.node.SendMessage(msg)
}

func (c *Client) WriteMessage(msg sock_msg.Message) (err error) {
	return c.node.WriteMessage(msg)
}

func (c *Client) WriteRaw(msg []byte) (int, error) {
	return c.node.WriteRaw(msg)
}

func (c *Client) Request(req sock_msg.Message) (sock_msg.Message, error) {
	return c.node.Request(req, c.timeout)
}

func (c *Client) Close(reason ...interface{}) {
	c.node.Stop(reason...)
	if c.onDisconnected != nil {
		c.onDisconnected(c.node)
	}
}

func (c *Client) IsClosed() bool {
	return c.node.closed.Val()
}

func (c *Client) StopReason() string {
	return c.node.StopReason()
}

func (c *Client) Endpoint() *Node { return c.node }

func (c *Client) Remote() string { return c.node.c.RemoteAddr().String() }

func (c *Client) State() NodeStat { return c.node.State() }
