package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"path"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/machinefi/w3bstream/pkg/modules/event"

	confapp "github.com/machinefi/w3bstream/pkg/depends/conf/app"
	conflog "github.com/machinefi/w3bstream/pkg/depends/conf/log"
	confmqtt "github.com/machinefi/w3bstream/pkg/depends/conf/mqtt"
	"github.com/machinefi/w3bstream/pkg/depends/protocol/eventpb"
)

var (
	broker = &confmqtt.Broker{}
	logger = conflog.Std()

	App *confapp.Ctx
)

func init() {
	App = confapp.New(
		confapp.WithName("mock-mqtt-client"),
		confapp.WithLogger(logger),
		confapp.WithRoot("."),
	)
	App.Conf(broker)

	flag.StringVar(&id, "id", "", "publish client id")
	flag.StringVar(&topic, "topic", "", "publish topic")
	flag.StringVar(&token, "token", "", "publish token")
	flag.StringVar(&data, "data", "", "payload data")
	flag.StringVar(&seq, "seq", "", "message sequence")
	flag.Parse()
}

var (
	id    string
	data  string
	topic string
	token string
	raw   []byte
	seq   string
	msg   *eventpb.Event
)

func init() {
	if seq == "" {
		seq = uuid.NewString()
	}

	var err error

	msg = &eventpb.Event{
		Header: &eventpb.Header{
			Token:   token,
			PubTime: time.Now().UTC().UnixMicro(),
			EventId: seq,
		},
		Payload: []byte(data),
	}

	raw, err = proto.Marshal(msg)
	if err != nil {
		panic(err)
	}
}

func main() {
	if id == "" {
		id = uuid.NewString()
	}
	c, err := broker.Client(id)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = c.WithTopic(topic).Publish(raw)
	if err != nil {
		fmt.Println(err)
		return
	}

	rspChannel := path.Join(topic, seq)

	err = c.WithTopic(rspChannel).Subscribe(func(cli mqtt.Client, msg mqtt.Message) {
		rsp := &event.EventRsp{}
		if err = json.Unmarshal(msg.Payload(), rsp); err != nil {
			fmt.Println(err)
		}
		ack, err := json.MarshalIndent(rsp, "", "  ")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(ack))
	})
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(time.Second * 3)
	_ = c.WithTopic(rspChannel).Unsubscribe()
}
