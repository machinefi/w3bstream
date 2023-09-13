package sock_test

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sock"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sock/sock_buf"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sock/sock_msg"
)

type parser struct {
}

func (p *parser) Marshal(buf sock_buf.Buffer, msg sock_msg.Message) error {
	dat, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = buf.Write(dat)
	if err != nil {
		return err
	}
	return nil
}

func (p *parser) Unmarshal(buf sock_buf.Buffer) (sock_msg.Message, error) {
	var (
		dat    = buf.Bytes()
		header = &Header{}
		msg    sock_msg.Message
	)

	err := json.Unmarshal(dat, header)
	if err != nil {
		return nil, err
	}
	switch header.Typ {
	case "greeting_req":
		msg = &GreetingReq{}
	case "greeting_rsp":
		msg = &GreetingRsp{}
	default:
		return nil, errors.New("unknown message")
	}

	err = json.Unmarshal(dat, msg)
	if err != nil {
		return nil, err
	}
	buf.Shift(len(dat))

	return msg, nil
}

var Parser = &parser{}

var (
	_ sock_msg.Message = (*GreetingReq)(nil)
	_ sock_msg.Message = (*GreetingRsp)(nil)
)

type ID string

func (v ID) String() string { return string(v) }

type Type string

func (v Type) String() string { return string(v) }

type Header struct {
	Id        string `json:"id"`
	Typ       string `json:"typ"`
	Timestamp int64  `json:"ts"`
}

type GreetingReq struct {
	Header
	Content string `json:"content"`
	Pid     int    `json:"pid"`
}

func NewGreetingReq(content string) *GreetingReq {
	return &GreetingReq{
		Header: Header{
			Id:        uuid.New().String(),
			Typ:       "greeting_req",
			Timestamp: time.Now().UnixNano(),
		},
		Content: content,
		Pid:     os.Getpid(),
	}
}

func (r *GreetingReq) ID() sock_msg.ID     { return ID(r.Id) }
func (r *GreetingReq) Type() sock_msg.Type { return Type(r.Typ) }
func (r *GreetingReq) GetTimestamp() int64 { return r.Timestamp }
func (r *GreetingReq) Renew()              { r.Id, r.Timestamp = uuid.New().String(), time.Now().UnixNano() }
func (r *GreetingReq) RenewWithID(id sock_msg.ID) {
	r.Id, r.Timestamp = id.String(), time.Now().UnixNano()
}

type GreetingRsp GreetingReq

func NewGreetingRsp(content string) *GreetingRsp {
	return &GreetingRsp{
		Header: Header{
			Id:        uuid.New().String(),
			Typ:       "greeting_rsp",
			Timestamp: time.Now().UnixNano(),
		},
		Content: content,
		Pid:     os.Getpid(),
	}
}

func (r *GreetingRsp) ID() sock_msg.ID     { return ID(r.Id) }
func (r *GreetingRsp) Type() sock_msg.Type { return Type(r.Typ) }
func (r *GreetingRsp) GetTimestamp() int64 { return r.Timestamp }
func (r *GreetingRsp) Renew()              { r.Id, r.Timestamp = uuid.New().String(), time.Now().UnixNano() }
func (r *GreetingRsp) RenewWithID(id sock_msg.ID) {
	r.Id, r.Timestamp = id.String(), time.Now().UnixNano()
}

const (
	MsgTypeGreetingReq uint32 = 1
	MsgTypeGreetingRsp uint32 = 2
)

type tcpParser struct {
}

var TCPParser = &tcpParser{}

func (r *tcpParser) Marshal(buf sock_buf.Buffer, msg sock_msg.Message) error {
	tmp := make([]byte, 8)
	dat, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	switch msg.Type().String() {
	case "greeting_req":
		binary.BigEndian.PutUint32(tmp[4:8], MsgTypeGreetingReq)
	case "greeting_rsp":
		binary.BigEndian.PutUint32(tmp[4:8], MsgTypeGreetingRsp)
	}

	binary.BigEndian.PutUint32(tmp[0:4], uint32(len(dat)))
	buf.Write(tmp)
	buf.Write(dat)
	return nil
}

func (r *tcpParser) Unmarshal(buf sock_buf.Buffer) (msg sock_msg.Message, err error) {
	var tmp []byte
	tmp, err = buf.Probe(8)
	if err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint32(tmp[0:4])
	typ := binary.BigEndian.Uint32(tmp[4:8])

	tmp, err = buf.Probe(int(8 + length))
	if err != nil {
		return nil, err
	}

	switch typ {
	case MsgTypeGreetingReq:
		msg = new(GreetingReq)
	case MsgTypeGreetingRsp:
		msg = new(GreetingRsp)
	default:
		return nil, errors.New("unknown message")
	}
	err = json.Unmarshal(tmp[8:], msg)
	if err != nil {
		return nil, err
	}
	buf.Shift(int(8 + length))
	return
}

var rsp = &GreetingRsp{
	Header:  Header{Typ: "greeting_rsp"},
	Content: "hello",
	Pid:     os.Getpid(),
}

var typ = Type("greeting_req")

func HandleGreeting(ev *sock.Event) {
	if ev == nil || ev.Payload() == nil {
		return
	}
	ep := ev.Node()
	pl := ev.Payload()
	rsp.RenewWithID(pl.ID())
	fmt.Printf("%v -> %s\n", ep.ID(), UnsafeJSONMarshalString(pl))
	fmt.Printf("%v <- %s\n", ep.ID(), UnsafeJSONMarshalString(rsp))
	ev.Response(rsp)
}
