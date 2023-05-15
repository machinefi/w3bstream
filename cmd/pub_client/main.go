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
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/retry"
	"github.com/pkg/errors"

	conflog "github.com/machinefi/w3bstream/pkg/depends/conf/log"
	confmqtt "github.com/machinefi/w3bstream/pkg/depends/conf/mqtt"
	"github.com/machinefi/w3bstream/pkg/depends/protocol/eventpb"
	"github.com/machinefi/w3bstream/pkg/modules/event"
)

var (
	broker *confmqtt.Broker
	logger = conflog.Std()

	// App *confapp.Ctx
)

func init() {
	// App = confapp.New(
	// 	confapp.WithName("mock-mqtt-client"),
	// 	confapp.WithLogger(logger),
	// 	confapp.WithRoot("."),
	// )
	// App.Conf(broker)

	flag.StringVar(&cid, "id", "", "publish client id")
	flag.StringVar(&topic, "topic", "", "publish topic")
	flag.StringVar(&token, "token", "", "publish token")
	flag.StringVar(&data, "data", "", "payload data")
	flag.StringVar(&seq, "seq", "", "message sequence")
	flag.StringVar(&host, "host", "localhost", "mqtt broker host")
	flag.IntVar(&port, "port", 1883, "mqtt broker port")
	flag.StringVar(&username, "username", "", "mqtt client username")
	flag.StringVar(&password, "password", "", "mqtt client password")
	flag.IntVar(&wait, "wait", 10, "mqtt wait ack timeout(seconds)")
	flag.Parse()
}

var (
	cid      string         // client id/pub id
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
	if cid == "" {
		cid = uuid.NewString()
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
		Timeout:   types.Duration(time.Second * 10),
		Keepalive: types.Duration(time.Minute),
		QoS:       confmqtt.QOS__ONCE,
	}

	broker.SetDefault()
	if err := broker.Init(); err != nil {
		panic(errors.Wrap(err, "init broker"))
	}

	var err error

	msg = &eventpb.Event{
		Header: &eventpb.Header{
			Token:   token,
			PubTime: time.Now().UTC().UnixMicro(),
			EventId: seq,
			PubId:   cid,
		},
		Payload: []byte(data),
	}

	raw, err = proto.Marshal(msg)
	if err != nil {
		panic(errors.Wrap(err, "build message"))
	}
}

func main() {
	c, err := broker.Client(cid)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = c.WithTopic(topic).Publish(raw)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(">>>> message published")

	rspChannel := path.Join(topic, cid)
	rspChan := make(chan interface{}, 0)

	err = c.WithTopic(rspChannel).Subscribe(func(cli mqtt.Client, msg mqtt.Message) {
		fmt.Println("<<<< message ack received")
		rsp := &event.EventRsp{}
		if err = json.Unmarshal(msg.Payload(), rsp); err != nil {
			fmt.Println(err)
		}
		ack, err := json.MarshalIndent(rsp, "", "  ")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(ack))
		rspChan <- 0
	})
	if err != nil {
		fmt.Println(err)
	}
	select {
	case <-rspChan:
	case <-time.After(time.Second * time.Duration(wait)):
		fmt.Println("**** message ack timeout")
	}
	_ = c.WithTopic(rspChannel).Unsubscribe()
}
