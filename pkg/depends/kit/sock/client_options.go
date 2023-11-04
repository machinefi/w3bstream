package sock

import (
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sock/sock_msg"
)

type ClientOption struct {
	remote          string          // remote addr
	protocol        ProtocolType    // protocol type
	timeout         time.Duration   // timeout rea/write timeout
	parser          sock_msg.Parser // payload data parser for data conversation
	readBufferCap   int             // user space read buffer capacity
	writeBufferCap  int             // user space write buffer capacity
	readBufferSize  int             // socket option for receive buffer size
	writeBufferSize int             // socket option for write buffer size
	nodeID          string          // client identifier
	workerPoolSize  int             // limited concurrent worker
	onConnected     func(*Node)     // udp ignored
	onDisconnected  func(*Node)     //
	routes          *Routes         // bind message type and handler
	handler         Handler         // common handler to invoke business API
	debug           bool
	output          func(interface{})
}

type ClientOptionSetter func(*ClientOption)

func ClientOptionRemote(v string) ClientOptionSetter {
	return func(o *ClientOption) {
		o.remote = v
	}
}

func ClientOptionProtocol(v ProtocolType) ClientOptionSetter {
	return func(o *ClientOption) {
		o.protocol = v
	}
}

func ClientOptionTimeout(v time.Duration) ClientOptionSetter {
	return func(o *ClientOption) {
		o.timeout = v
	}
}

func ClientOptionReadBufferCap(v int) ClientOptionSetter {
	return func(o *ClientOption) {
		o.readBufferCap = v
	}
}

func ClientOptionWriteBufferCap(v int) ClientOptionSetter {
	return func(o *ClientOption) {
		o.writeBufferCap = v
	}
}

func ClientOptionParser(v sock_msg.Parser) ClientOptionSetter {
	return func(o *ClientOption) {
		o.parser = v
	}
}

func ClientOptionNodeID(v string) ClientOptionSetter {
	return func(o *ClientOption) {
		o.nodeID = v
	}
}

func ClientOptionReadBufferSize(v int) ClientOptionSetter {
	return func(o *ClientOption) {
		o.readBufferSize = v
	}
}

func ClientOptionWriteBufferSize(v int) ClientOptionSetter {
	return func(o *ClientOption) {
		o.writeBufferSize = v
	}
}

func ClientOptionWorkerPoolSize(v int) ClientOptionSetter {
	return func(o *ClientOption) {
		o.workerPoolSize = v
	}
}

func ClientOptionOnConnected(f func(*Node)) ClientOptionSetter {
	return func(o *ClientOption) {
		o.onConnected = f
	}
}

func ClientOptionOnDisconnected(f func(*Node)) ClientOptionSetter {
	return func(o *ClientOption) {
		o.onDisconnected = f
	}
}

func ClientOptionRoute(t sock_msg.Type, h ...Handler) ClientOptionSetter {
	return func(o *ClientOption) {
		if o.routes == nil {
			o.routes = NewRoutes()
		}
		o.routes.Register(t, h...)
	}
}

func ClientOptionDebugMode(output func(interface{})) ClientOptionSetter {
	return func(o *ClientOption) {
		o.debug = true
		o.output = output
	}
}

func ClientOptionHandler(h Handler) ClientOptionSetter {
	return func(o *ClientOption) {
		o.handler = h
	}
}
