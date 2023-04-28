package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"

	confapp "github.com/machinefi/w3bstream/pkg/depends/conf/app"
	conflog "github.com/machinefi/w3bstream/pkg/depends/conf/log"
	confmqtt "github.com/machinefi/w3bstream/pkg/depends/conf/mqtt"
	"github.com/machinefi/w3bstream/pkg/depends/protocol/eventpb"
)

var (
	mqtt   = &confmqtt.Broker{}
	logger = conflog.Std()

	App *confapp.Ctx
)

func init() {
	App = confapp.New(
		confapp.WithName("mock-mqtt-client"),
		confapp.WithLogger(logger),
		confapp.WithRoot("."),
	)
	App.Conf(mqtt)

	flag.StringVar(&id, "id", "", "publish client id")
	flag.StringVar(&topic, "topic", "", "publish topic")
	flag.StringVar(&token, "token", "", "publish token")
	flag.StringVar(&data, "data", "", "payload data")
	flag.Parse()
}

var (
	id      string
	data    string
	topic   string
	token   string
	payload string
	raw     []byte
)

func init() {
	var err error
	raw, err = proto.Marshal(&eventpb.Event{
		Header: &eventpb.Header{
			Token:   token,
			PubTime: time.Now().UTC().UnixMicro(),
			EventId: uuid.NewString(),
		},
		Payload: []byte(data),
	})

	if err != nil {
		panic(err)
	}
}

func main() {
	if id == "" {
		id = uuid.NewString()
	}
	c, err := mqtt.Client(id)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = c.WithTopic(topic).Publish(raw)
	if err != nil {
		fmt.Println(err)
		return
	}
}
