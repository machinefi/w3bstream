package sock

import "github.com/machinefi/w3bstream/pkg/depends/kit/sock/sock_msg"

type Event struct {
	node    *Node
	payload sock_msg.Message
}

func NewEvent(node *Node, pld sock_msg.Message) *Event {
	return &Event{node, pld}
}

func (ev *Event) Node() *Node { return ev.node }

func (ev *Event) Payload() sock_msg.Message { return ev.payload }

func (ev *Event) Response(rsp sock_msg.Message) error {
	return ev.node.WriteMessage(rsp)
}

func (ev *Event) Send(msg sock_msg.Message) error {
	return ev.node.SendMessage(msg)
}
