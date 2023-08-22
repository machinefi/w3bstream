package types

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"strconv"
	"sync/atomic"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/x/ptrx"
)

type Publisher struct {
	DeviceMN string `env:""`
	Token    string `env:""`
}

func (c *Publisher) IsZero() bool {
	return c.DeviceMN == "" && c.Token == ""
}

type Channel struct {
	Publisher   `env:""`
	Topic       string        `env:""`
	EventType   string        `env:""`
	PubInterval time.Duration `env:""`
}

func (c *Channel) SetDefault() {
	if c.EventType == "" {
		c.EventType = "DEFAULT"
	}
	if c.PubInterval == 0 {
		c.PubInterval = time.Second
	}
}

func (c *Channel) IsZero() bool {
	return c.Publisher.IsZero() || c.Topic == ""
}

func (c *Channel) Subscribe(ctx context.Context) error {
	broker := MustMqttBrokerFromContext(ctx)
	topic := path.Join("ack", c.DeviceMN)

	client, err := broker.Client(uuid.NewString())
	if err != nil {
		return err
	}
	var (
		succeeded = ptrx.Ptr(int64(0))
		total     = ptrx.Ptr(int64(0))
	)

	err = client.WithTopic(topic).WithQoS(broker.QoS).Subscribe(func(_ mqtt.Client, msg mqtt.Message) {
		var (
			// ts1     = time.Now().UTC().UnixMilli() // ts1 timestamp client sub ack
			statOk  = int64(0)
			statAll = atomic.AddInt64(total, 1)
		)
		body := &struct {
			EventID   string `json:"eventID"`
			Timestamp int64  `json:"timestamp"` // ts server pub ack
		}{}
		if err = json.Unmarshal(msg.Payload(), body); err == nil {
			statOk = atomic.AddInt64(succeeded, 1)
		}
		fmt.Printf("%s <<< %s %d/%d\n", c.DeviceMN, body.EventID, statOk, statAll)
	})
	if err != nil {
		fmt.Printf("subscrib %s failed [err: %v]\n", topic, err)
		return err
	}
	fmt.Printf("start subscribing %s\n", topic)
	return nil
}

func (c *Channel) StartPublish(ctx context.Context) {
	broker := MustMqttBrokerFromContext(ctx)

	client, err := broker.Client(uuid.NewString())
	if err != nil {
		panic(err)
	}

	seq := 1

	for {
		eventID := fmt.Sprintf("%s_%010d", c.DeviceMN, seq)
		ts := strconv.FormatInt(time.Now().UTC().UnixMilli(), 10)
		topic := fmt.Sprintf("%s/push/%s/%s/id=%s&ts=%s", c.Topic, c.Token, c.EventType, eventID, ts)
		err = client.WithQoS(broker.QoS).WithTopic(topic).Publish(ts)
		if err != nil {
			fmt.Printf("%v\n", errors.Wrap(err, "publish message"))
			goto Next
		}
		fmt.Printf("%s >>> seq: %s\n", c.DeviceMN, eventID)
	Next:
		seq++
		time.Sleep(c.PubInterval)
	}
}
