package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	conflog "github.com/machinefi/w3bstream/pkg/depends/conf/log"
	confmqtt "github.com/machinefi/w3bstream/pkg/depends/conf/mqtt"
	"github.com/machinefi/w3bstream/pkg/depends/protocol/eventpb"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/retry"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/timer"
	"github.com/machinefi/w3bstream/pkg/modules/event"
)

var (
	broker *confmqtt.Broker
	logger = conflog.Std()
)

func init() {
	flag.StringVar(&mn, "mn", "", "publisher(device) mn")
	flag.StringVar(&topic, "topic", "", "publish topic")
	flag.StringVar(&token, "token", "", "publish token")
	flag.StringVar(&data, "data", "", "payload data, read file pls use '@PATH'")
	flag.StringVar(&seq, "seq", "", "message sequence")
	flag.StringVar(&host, "host", "localhost", "mqtt broker host")
	flag.IntVar(&port, "port", 1883, "mqtt broker port")
	flag.StringVar(&username, "username", "", "mqtt client username")
	flag.StringVar(&password, "password", "", "mqtt client password")
	flag.IntVar(&wait, "wait", 10, "mqtt wait ack timeout(seconds)")
	flag.Parse()
}

var (
	mn       string         // publisher mn
	data     string         // message payload
	topic    string         // mqtt topic
	token    string         // publisher token
	host     string         // mqtt broker host
	port     int            // mqtt broker port
	username string         // mqtt client username
	password string         // mqtt client password
	wait     int            // mqtt wait ack timeout
	seq      string         // message sequence
	raw      []byte         // mqtt message
	msg      *eventpb.Event // mqtt message (protobuf)
)

func init() {
	if seq == "" {
		seq = uuid.NewString()
	}
	if mn == "" {
		mn = uuid.NewString()
	}
	if host == "" {
		host = "localhost"
	}
	if port == 0 {
		port = 1883
	}

	broker = &confmqtt.Broker{
		Server: types.Endpoint{
			Scheme:   "mqtt",
			Hostname: host,
			Port:     uint16(port),
			Username: username,
			Password: types.Password(password),
		},
		Retry:     *retry.Default,
		Timeout:   types.Duration(time.Second * time.Duration(wait)),
		Keepalive: types.Duration(time.Second * time.Duration(wait)),
		QoS:       confmqtt.QOS__ONLY_ONCE,
	}

	broker.SetDefault()
	if err := broker.Init(); err != nil {
		panic(errors.Wrap(err, "init broker"))
	}

	var err error

	pl := []byte(data)
	if len(data) > 0 && data[0] == '@' {
		pl, err = os.ReadFile(data[1:])
		if err != nil {
			panic(errors.Wrap(err, "read file: "+data[1:]))
		}
	}

	msg = &eventpb.Event{
		Header: &eventpb.Header{
			Token:   token,
			PubTime: time.Now().UTC().UnixMicro(),
			EventId: seq,
			PubId:   mn,
		},
		Payload: pl,
	}

	raw, err = proto.Marshal(msg)
	if err != nil {
		panic(errors.Wrap(err, "build message"))
	}
}

func main() {
	c, err := broker.Client(mn)
	if err != nil {
		fmt.Println(err)
		return
	}
	cost := timer.Start()
	err = c.WithTopic(topic).Publish(raw)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(">>>> message published")

	ack := path.Join("ack", mn)
	fmt.Printf("receive ack from %s\n", ack)
	sig := make(chan interface{})

	err = c.WithTopic(ack).Subscribe(func(cli mqtt.Client, msg mqtt.Message) {
		fmt.Println("<<<< message ack received")
		rsp := &event.EventRsp{}
		if err = json.Unmarshal(msg.Payload(), rsp); err != nil {
			fmt.Println(err)
		}
		content, _ := json.MarshalIndent(rsp, "", "  ")
		if err != nil {
			fmt.Printf("unmarshal failed %v\n", err)
		} else {
			fmt.Println(string(content))
		}
		sig <- 0
	})
	if err != nil {
		fmt.Println(err)
	}
	select {
	case <-sig:
	case <-time.After(time.Second * time.Duration(wait)):
		fmt.Println("**** message ack timeout")
	}
	fmt.Printf("cost: %dms", cost().Milliseconds())
	_ = c.WithTopic(ack).Unsubscribe()
}
