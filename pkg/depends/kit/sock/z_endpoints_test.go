package sock_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sock"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sock/sock_msg"
)

func UnsafeJSONMarshalString(v interface{}) string {
	content, _ := json.Marshal(v)
	return string(content)
}

func RunTCPClient(ctx context.Context) {
	var req = NewGreetingReq("hello")
	var cli, err = sock.NewClient(
		sock.ClientOptionParser(TCPParser),
		sock.ClientOptionRemote("localhost:10086"),
		sock.ClientOptionProtocol(sock.ProtocolTCP),
		sock.ClientOptionOnConnected(func(n *sock.Node) {
			req.Content = "first greeting req"
			for {
				req.Renew()
				rsp, err := n.Request(req, 2*time.Second)
				if err != nil {
					fmt.Printf("cli: OnConnected error:%v\n", err)
					if sock.IsTimeoutError(err) {
						continue
					}
					if sock.IsNodeClosedError(err) {
						break
					}
				}
				msg, ok := rsp.(*GreetingRsp)
				if ok && msg.Content == "first greeting rsp" {
					fmt.Println("OnConnected done")
					break
				}
			}
		}),
		sock.ClientOptionRoute(rsp.Type(), func(ev *sock.Event) {
			fmt.Printf("route -> %s\n", UnsafeJSONMarshalString(ev.Payload()))
		}),
	)

	if err != nil {
		panic(err)
	}

	req.Content = "hello"
	for {
		select {
		case <-ctx.Done():
			cli.Close()
			return
		default:
		}
		fmt.Printf("-> %s\n", UnsafeJSONMarshalString(req))
		rsp, err := cli.Request(req)
		if err != nil {
			if sock.IsNodeClosedError(err) {
				break
			}
			fmt.Println(err)
			time.Sleep(3 * time.Second)
			continue
		}
		fmt.Printf("<- %s, %d\n", UnsafeJSONMarshalString(rsp),
			(rsp.(sock_msg.WithTimestamp).GetTimestamp()-req.GetTimestamp())/1e6)
		time.Sleep(250 * time.Millisecond)
		req.Renew()
	}
}

func RunTCPServer(ctx context.Context) {
	var addr = ":10086"
	var srv, err = sock.NewServer(
		sock.ServerOptionConnCap(10),
		sock.ServerOptionListenAddr(addr),
		sock.ServerOptionProtocol(sock.ProtocolTCP),
		sock.ServerOptionParser(TCPParser),
		sock.ServerOptionRoute(typ, HandleGreeting),
	)
	sock.ServerOptionOnConnected(func(n *sock.Node) {
		for {
			msg, err := n.ReadMessage(time.Second)
			if err != nil {
				fmt.Printf("srv: OnConnected [error:%v]\n", err)
				if sock.IsTimeoutError(err) {
					continue
				}
				if sock.IsNodeClosedError(err) {
					break
				}
			}
			raw, ok := msg.(*GreetingReq)
			if ok && raw.Content == "first greeting req" {
				err = n.WriteMessage(&GreetingRsp{
					Header: Header{
						Id:        raw.Id,
						Typ:       "greeting_rsp",
						Timestamp: time.Now().UnixNano(),
					},
					Content: "first greeting rsp",
					Pid:     os.Getegid(),
				})
				if err != nil {
					fmt.Printf("srv: OnConnected [error:%v]\n", err)
					break
				}
				fmt.Println("srv: OnConnected done")
				break
			}
			continue
		}
	})

	fmt.Println("server started: " + addr)

	if err != nil {
		panic(err)
	}
	srv.Serve(ctx)
}

func RunUDPClient(ctx context.Context) {
	var req = NewGreetingReq("first greeting req")
	var cli, err = sock.NewClient(
		sock.ClientOptionParser(Parser),
		sock.ClientOptionRemote("localhost:10010"),
		sock.ClientOptionProtocol(sock.ProtocolUDP),
		sock.ClientOptionNodeID("udp_client"),
		sock.ClientOptionTimeout(time.Second),
		sock.ClientOptionOnConnected(func(n *sock.Node) {
			req.Content = "first greeting req"
			for {
				req.Renew()
				rsp, err := n.Request(req, 2*time.Second)
				if err != nil {
					fmt.Printf("cli: OnConnected error:%v\n", err)
					if sock.IsTimeoutError(err) {
						continue
					}
					if sock.IsNodeClosedError(err) {
						break
					}
				}
				msg, ok := rsp.(*GreetingRsp)
				if ok && msg.Content == "first greeting rsp" {
					fmt.Println("OnConnected done")
					break
				}
			}
		}),
		sock.ClientOptionRoute(rsp.Type(), func(ev *sock.Event) {
			fmt.Printf("route -> %s\n", UnsafeJSONMarshalString(ev.Payload()))
		}),
	)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case <-ctx.Done():
			cli.Close()
			return
		default:
		}
		if cli.IsClosed() {
			break
		}
		fmt.Printf("-> %s\n", UnsafeJSONMarshalString(req))
		rsp, err := cli.Request(req)
		if err != nil {
			if cli.IsClosed() {
				break
			}
			fmt.Println(err)
			time.Sleep(3 * time.Second)
			continue
		}
		fmt.Printf("<- %s, %d\n", UnsafeJSONMarshalString(rsp),
			rsp.(sock_msg.WithTimestamp).GetTimestamp()-req.GetTimestamp())
		time.Sleep(time.Millisecond * 250)
		req.Renew()
	}
	fmt.Println("client closed")
}

func RunUDPServer(ctx context.Context) {
	var addr = ":10010"
	var srv, err = sock.NewServer(
		sock.ServerOptionConnCap(10),
		sock.ServerOptionListenAddr(addr),
		sock.ServerOptionProtocol(sock.ProtocolUDP),
		sock.ServerOptionParser(Parser),
		sock.ServerOptionRoute(typ, HandleGreeting),
		sock.ServerOptionOnConnected(func(n *sock.Node) {
			for {
				msg, err := n.ReadMessage(time.Second)
				if err != nil {
					fmt.Printf("srv: OnConnected [error:%v]\n", err)
					if sock.IsTimeoutError(err) {
						continue
					}
					if sock.IsNodeClosedError(err) {
						break
					}
				}
				raw, ok := msg.(*GreetingReq)
				if ok && raw.Content == "first greeting req" {
					err = n.WriteMessage(&GreetingRsp{
						Header: Header{
							Id:        raw.Id,
							Typ:       "greeting_rsp",
							Timestamp: time.Now().UnixNano(),
						},
						Content: "first greeting rsp",
						Pid:     os.Getegid(),
					})
					if err != nil {
						fmt.Printf("srv: OnConnected [error:%v]\n", err)
						break
					}
					fmt.Println("srv: OnConnected done")
					break
				}
				continue
			}
		}),
	)

	fmt.Println("server started: " + addr)

	if err != nil {
		panic(err)
	}

	srv.Serve(ctx)
}

func RunUDSClient(ctx context.Context) {
	var cli, err = sock.NewClient(
		sock.ClientOptionParser(Parser),
		sock.ClientOptionRemote("/tmp/ipc.sock"),
		sock.ClientOptionProtocol(sock.ProtocolUnix),
	)
	if err != nil {
		panic(err)
	}

	req := NewGreetingReq("hello")
	for {
		select {
		case <-ctx.Done():
			cli.Close()
			return
		default:
		}

		fmt.Printf("-> %s\n", UnsafeJSONMarshalString(req))
		rsp, err := cli.Request(req)
		if err != nil {
			fmt.Println(err)
			time.Sleep(3 * time.Second)
			continue
		}
		fmt.Printf("<- %s, %d\n", UnsafeJSONMarshalString(rsp),
			rsp.(sock_msg.WithTimestamp).GetTimestamp()-req.GetTimestamp())
		time.Sleep(time.Second)
		req.Renew()
	}
}

func RunUDSServer(ctx context.Context) {
	var addr = "/tmp/ipc.sock"
	var srv, err = sock.NewServer(
		sock.ServerOptionConnCap(10),
		sock.ServerOptionListenAddr(addr),
		sock.ServerOptionProtocol(sock.ProtocolUnix),
		sock.ServerOptionParser(Parser),
		sock.ServerOptionRoute(typ, HandleGreeting))

	fmt.Println("server started: " + addr)

	if err != nil {
		panic(err)
	}
	srv.Serve(ctx)
}

// test tcp
func TestTCP(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	go RunTCPServer(ctx)
	time.Sleep(time.Second)
	go RunTCPClient(ctx)

	time.Sleep(time.Second * 5)
	cancel()
}

// test udp
func TestUDP(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	go RunUDPServer(ctx)
	time.Sleep(time.Second)
	go RunUDPClient(ctx)

	time.Sleep(time.Second * 5)
	cancel()
}

// test unix domain socket
func TestUDS(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	go RunUDSServer(ctx)
	time.Sleep(time.Second)
	go RunUDSClient(ctx)

	time.Sleep(time.Second * 5)
	cancel()
}
