package sock

import (
	"fmt"
	"sync"
)

type clients struct {
	*sync.RWMutex
	connections map[string]*Node
	events      chan *Event
}

func NewClientManager() *clients {
	return &clients{
		RWMutex:     &sync.RWMutex{},
		connections: make(map[string]*Node),
		events:      make(chan *Event, 2048),
	}
}

func (cs *clients) New(n *Node, onDisconnected ...func(*Node)) {
	cs.Lock()
	defer cs.Unlock()

	cs.connections[n.ID()] = n
	fmt.Printf("[sock.clients] new node added `%s`", n.id)

	go func(n *Node) {
		defer cs.Remove(n.ID())
		defer func() {
			for _, f := range onDisconnected {
				f(n)
			}
		}()
		for {
			if n.closed.Val() {
				return
			}
			msg, _ := n.ReadMessage(0)
			if msg == nil {
				break
			}
			cs.events <- &Event{n, msg}
		}
	}(n)
}

func (cs *clients) Remove(id string, reason ...interface{}) {
	cs.Lock()
	defer cs.Unlock()

	n := cs.connections[id]
	if n != nil {
		fmt.Printf("[sock.clients] node leave `%s`", n.id)
		n.Stop(reason...)
	}
	delete(cs.connections, id)
}

func (cs *clients) Get(id string) *Node {
	cs.RLock()
	defer cs.RUnlock()
	ret, ok := cs.connections[id]
	if ok {
		return ret
	}
	return nil
}

func (cs *clients) Stats() (int, int) {
	cs.RLock()
	defer cs.RUnlock()

	connections, events := len(cs.connections), len(cs.events)
	fmt.Printf("[sock.clients] connections: %d, events: %d", connections, events)
	return connections, events
}

func (cs *clients) Reset() {
	cs.Lock()
	defer cs.Unlock()

	for id, n := range cs.connections {
		if n != nil {
			n.Stop("server reset clients")
		}
		delete(cs.connections, id)
	}
}
