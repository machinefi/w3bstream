package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"sync"
	"sync/atomic"
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
	flag.IntVar(&conc, "conc", 10, "")
	flag.IntVar(&count, "count", 10, "")
	flag.Parse()
}

var (
	mn       string         // publisher mn
	conc     int            // conc concurrency
	count    int            // count message count
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
		QoS:       confmqtt.QOS__ONCE,
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

func cli(mn string) (*confmqtt.Client, error) {
	c, err := broker.Client(mn)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func sub(c *confmqtt.Client, topic string, succeed *int64) error {
	return c.WithTopic(topic).Subscribe(func(cli mqtt.Client, msg mqtt.Message) {
		fmt.Printf("%s <<< message subscribed\n", mn)
		atomic.AddInt64(succeed, 1)
	})
}

func pub(c *confmqtt.Client, mn, topic string, count int) (int64, time.Duration) {
	succeed := int64(0)
	wg := &sync.WaitGroup{}

	cost := timer.Start()
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			err := c.WithTopic(topic).Publish(raw)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("%s >>> message published\n", mn)

			wg.Done()
		}()
	}
	wg.Wait()
	return succeed, cost()
}

func main() {
	wg := &sync.WaitGroup{}

	succeed := int64(0)
	wg.Add(conc)
	for i := 0; i < conc; i++ {
		c, err := cli(mn)
		if err != nil {
			fmt.Println(err)
			continue
		}

		pubTopic := topic
		subTopic := path.Join("ack", mn)

		if err := sub(c, subTopic, &succeed); err != nil {
			fmt.Println(err)
			continue
		}

		go func(mn string) {
			n, _ := pub(c, mn, pubTopic, count)
			atomic.AddInt64(&succeed, n)
			wg.Done()
		}(mn)
	}
	wg.Wait()

	time.Sleep(time.Duration(wait) * time.Second)

	fmt.Printf("publishing: %d/%d\n", succeed, conc*count)
}
