package sock

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sche"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sche/safe"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sock/sock_buf"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sock/sock_msg"
)

type Node struct {
	c         net.Conn
	addr      net.Addr
	rb, wb    sock_buf.Buffer
	rq, sq    chan sock_msg.Message
	binder    *Binder
	routes    *Routes // routes  register sock_msg.Type and Handler high priority
	handler   Handler // handler register global message handle func low priority
	parser    sock_msg.Parser
	worker    sche.WorkersScheduler
	id        string
	ctx       context.Context
	cancel    context.CancelFunc
	onceClose *sync.Once
	writerMtx *sync.Mutex
	mgr       *clients
	protocol  ProtocolType
	closed    *safe.Bool
	reason    string
	debug     bool
	output    func(interface{})
}

func NewNode() *Node {
	ret := &Node{
		onceClose: &sync.Once{},
		writerMtx: &sync.Mutex{},
		binder:    NewBinder(),
		closed:    safe.NewBool(),
	}
	ret.ctx, ret.cancel = context.WithCancel(context.Background())
	return ret
}

func (n *Node) ID() string { return n.id }

func (n *Node) Context() context.Context { return n.ctx }

func (n *Node) SetQueues(r, s chan sock_msg.Message) { n.rq, n.sq = r, s }

func (n *Node) SetParser(parser sock_msg.Parser) { n.parser = parser }

func (n *Node) SetBuffers(r, w sock_buf.Buffer) { n.rb, n.wb = r, w }

func (n *Node) WithContext(ctx context.Context) *Node {
	n.ctx, n.cancel = context.WithCancel(ctx)
	return n
}

// ReadMessage read message from receive channel
func (n *Node) ReadMessage(d time.Duration) (msg sock_msg.Message, err error) {
	if n.closed.Val() {
		return nil, ENodeClosed
	}
	if d == 0 {
		return <-n.rq, nil
	}
	select {
	case msg = <-n.rq:
		return msg, nil
	case <-time.After(d):
		return nil, ENodeTimeout
	}
}

// WriteMessage write message to send buffer
func (n *Node) WriteMessage(msg sock_msg.Message) (err error) {
	if n.closed.Val() {
		return ENodeClosed
	}
	n.writerMtx.Lock()
	defer n.writerMtx.Unlock()
	n.wb.Reset()
	err = n.parser.Marshal(n.wb, msg)
	if err != nil {
		return ENodeMarshal.WithError(err)
	}
	var dat = n.wb.Bytes()
	_, err = n.write(dat)
	if err != nil {
		err = ENodeWrite.WithError(err)
		n.Stop(err)
	}
	return err
}

// WriteRaw write binary data to send buffer
func (n *Node) WriteRaw(msg []byte) (int, error) {
	if n.closed.Val() {
		return 0, ENodeClosed
	}
	n.writerMtx.Lock()
	defer n.writerMtx.Unlock()
	return n.write(msg)
}

// SendMessage push message to send queue
func (n *Node) SendMessage(msg sock_msg.Message) (err error) {
	if n.closed.Val() {
		return ENodeClosed
	}
	defer func() {
		e, ok := recover().(error)
		if ok && e != nil {
			err = e
		}
	}()
	n.sq <- msg
	return nil
}

// Request request message to peer until timeout or responded
func (n *Node) Request(req sock_msg.Message, d time.Duration) (sock_msg.Message, error) {
	if n.closed.Val() {
		return nil, ENodeClosed
	}
	if n.binder.New(req.ID()) != nil {
		return nil, EMessageIdRepeated
	}
	if err := n.WriteMessage(req); err != nil {
		n.binder.Remove(req.ID())
		return nil, err
	}
	return n.binder.Wait(req.ID(), d)
}

// Stop stop Node running
func (n *Node) Stop(reason ...interface{}) {
	n.onceClose.Do(func() {
		if n.cancel != nil {
			n.cancel()
		}
		addr := ""
		if n.c != nil && n.addr == nil {
			addr = n.c.RemoteAddr().String()
			_ = n.c.Close()
		}
		if n.addr != nil {
			addr = n.addr.String()
		}
		if n.sq != nil {
			close(n.sq)
		}
		if n.rq != nil {
			close(n.rq)
		}
		if n.wb != nil {
			n.wb.Reset()
		}
		if n.rb != nil {
			n.rb.Reset()
		}
		if n.binder != nil {
			n.binder.Reset()
		}
		n.closed.Set(true)
		n.reason = fmt.Sprint(reason...)
		fmt.Printf("%s cloesd: %v remote: %s\n", n.id, n.reason, addr)
	})
}

// StopReason return stop reason
func (n *Node) StopReason() string { return n.reason }

// Run start Node writing and reading
func (n *Node) Run() {
	if n.protocol != ProtocolUDP {
		go n.recv()
	}
	go n.send()
}

func (n *Node) RunUDP() { go n.send() }

func (n *Node) send() {
	var (
		msg sock_msg.Message
		err error
	)

	defer func() {
		fmt.Printf("[qsock.send] error %v\n", err)
		n.Stop(err)
	}()

	for {
		select {
		case msg = <-n.sq:
		case <-n.ctx.Done():
			err = ENodeContextDone
			return
		}
		if msg == nil {
			err = ENodeMessage
			return
		}
		if n.debug && n.output != nil {
			n.output(msg)
		}
		if err = n.WriteMessage(msg); err != nil {
			e := errors.Unwrap(err)
			if errors.Is(e, ENodeWrite) {
				return
			}
		}
	}
}

func (n *Node) recv() {
	var (
		msg  sock_msg.Message
		err  error
		dat  = make([]byte, DefaultBufferCapacity)
		size int
	)

	defer func() {
		n.Stop(err)
	}()

	for {
		select {
		case <-n.ctx.Done():
			err = n.ctx.Err()
			return
		default:
			size, err = n.read(dat)
			if err != nil {
				err = ENodeRead.WithError(err)
				return
			}
			if size == 0 {
				err = ENodeRead.WithMessage("read size=0")
				return
			}

			err = n.rb.ResetAndWrite(dat[0:size])
			dat = dat[0:cap(dat)]
			for {
				msg, err = n.parser.Unmarshal(n.rb)
				if err != nil {
					if errors.Is(err, sock_buf.EStreamBufferDataLack) {
						break
					}
					return
				}
				if n.binder.Push(msg) {
					goto check
				}
				if handlers := n.routes.Handlers(msg.Type()); len(handlers) > 0 {
					for _, handler := range handlers {
						n.worker.Add(HandlerFunc(handler, &Event{n, msg}))
					}
					goto check
				}
				if n.handler != nil {
					n.worker.Add(HandlerFunc(n.handler, &Event{n, msg}))
					goto check
				}
				n.rq <- msg
			check:
				if n.protocol == ProtocolTCP {
					continue
				}
				break
			}
		}
	}
}

func (n *Node) read(dat []byte) (size int, err error) {
	if n.addr == nil {
		size, err = n.c.Read(dat)
	} else {
		size, _, err = n.c.(*net.UDPConn).ReadFromUDP(dat)
	}
	return
}

func (n *Node) write(dat []byte) (size int, err error) {
	if n.addr != nil {
		size, err = n.c.(*net.UDPConn).WriteTo(dat, n.addr)
		return
	} else {
		written := 0
		for size < len(dat) {
			written, err = n.c.Write(dat[size:])
			if err != nil {
				return
			}
			size += written
		}
	}
	return
}

type NodeStat struct {
	RBufLen int `json:"rb_len"`
	RBufCap int `json:"rb_cap"`
	WBufLen int `json:"wb_len"`
	WBufCap int `json:"wb_cap"`
	RQLen   int `json:"rq_len"`
	SQLen   int `json:"sq_len"`
}

func (n *Node) State() NodeStat {
	return NodeStat{
		n.rb.Len(),
		n.rb.Cap(),
		n.wb.Len(),
		n.wb.Cap(),
		len(n.rq),
		len(n.sq),
	}
}
