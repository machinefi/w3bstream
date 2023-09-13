package sock

import (
	"net"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sock/sock_msg"
)

type ServerOption struct {
	listenAddr      string
	protocol        ProtocolType
	timeout         time.Duration
	parser          sock_msg.Parser
	connCap         int
	readBufferCap   int
	writeBufferCap  int
	readBufferSize  int
	writeBufferSize int
	onConnected     []func(*Node) // udp ignored
	onDisconnected  []func(*Node)
	workerPoolSize  int
	routes          *Routes
	handler         Handler
	debug           bool
	output          func(interface{})
}

type ServerOptionSetter func(*ServerOption)

func ServerOptionListenAddr(v string) ServerOptionSetter {
	return func(o *ServerOption) {
		o.listenAddr = v
	}
}

func ServerOptionProtocol(v ProtocolType) ServerOptionSetter {
	return func(o *ServerOption) {
		o.protocol = v
	}
}

func ServerOptionTimeout(v time.Duration) ServerOptionSetter {
	return func(o *ServerOption) {
		o.timeout = v
	}
}

func ServerOptionParser(v sock_msg.Parser) ServerOptionSetter {
	return func(o *ServerOption) {
		o.parser = v
	}
}

func ServerOptionConnCap(v int) ServerOptionSetter {
	return func(o *ServerOption) {
		o.connCap = v
	}
}

func ServerOptionReadBufferCap(v int) ServerOptionSetter {
	return func(o *ServerOption) {
		o.readBufferCap = v
	}
}

func ServerOptionWriteBufferCap(v int) ServerOptionSetter {
	return func(o *ServerOption) {
		o.writeBufferCap = v
	}
}

func ServerOptionOnConnected(f func(*Node)) ServerOptionSetter {
	return func(o *ServerOption) {
		o.onConnected = append(o.onConnected, f)
	}
}

func ServerOptionOnDisconnected(f func(*Node)) ServerOptionSetter {
	return func(o *ServerOption) {
		o.onDisconnected = append(o.onDisconnected, f)
	}
}

func ServerOptionReadBufferSize(v int) ServerOptionSetter {
	return func(o *ServerOption) {
		o.readBufferSize = v
	}
}

func ServerOptionWriteBufferSize(v int) ServerOptionSetter {
	return func(o *ServerOption) {
		o.writeBufferSize = v
	}
}

func ServerOptionWorkerPoolSize(v int) ServerOptionSetter {
	return func(o *ServerOption) {
		o.workerPoolSize = v
	}
}

func ServerOptionRoute(t sock_msg.Type, h ...Handler) ServerOptionSetter {
	return func(o *ServerOption) {
		if o.routes == nil {
			o.routes = NewRoutes()
		}
		o.routes.Register(t, h...)
	}
}

func ServerOptionHandler(h Handler) ServerOptionSetter {
	return func(o *ServerOption) {
		o.handler = h
	}
}

func ServerOptionDebugMode(output func(interface{})) ServerOptionSetter {
	return func(o *ServerOption) {
		o.debug = true
		o.output = output
	}
}

const (
	DefaultBufferCapacity = 1024 * 1024
)

func SetConnOption(c net.Conn, r, w int) (err error) {
	switch v := c.(type) {
	case *net.TCPConn:
		if err = v.SetReadBuffer(r); err != nil {
			return
		}
		if err = v.SetWriteBuffer(w); err != nil {
			return
		}
		if err = v.SetKeepAlive(true); err != nil {
			return
		}
		if err = v.SetNoDelay(true); err != nil {
			return
		}
	case *net.UDPConn:
		if err = v.SetReadBuffer(r); err != nil {
			return
		}
		if err = v.SetWriteBuffer(w); err != nil {
			return
		}
	case *net.UnixConn:
		if err = v.SetReadBuffer(r); err != nil {
			return
		}
		if err = v.SetWriteBuffer(w); err != nil {
			return
		}
	}
	return nil
}
