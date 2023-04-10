package mqtt_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	. "github.com/machinefi/w3bstream/pkg/depends/conf/mqtt"
)

type PayloadBody struct {
	EventID      string
	PubTimestamp int64
	Message      string
}

func NewPayloadBody(msg string) *PayloadBody {
	return &PayloadBody{
		EventID:      uuid.New().String(),
		PubTimestamp: time.Now().UnixMilli(),
		Message:      msg,
	}
}

func UnsafeJsonMarshal(v interface{}) []byte {
	data, _ := json.Marshal(v)
	return data
}

var (
	topic  = "test_demo"
	broker = &Broker{}
)

func init() {
	server := types.Endpoint{}

	err := server.UnmarshalText([]byte("mqtt://broker.emqx.io:1883"))
	if err != nil {
		panic(err)
	}

	broker.SetDefault()
	err = broker.Init()
	if err != nil {
		panic(err)
	}
}

func TestBroker(t *testing.T) {
	cpub, err := broker.Client("pub")
	NewWithT(t).Expect(err).To(BeNil())
	cpub.WithTopic(topic)

	csub, err := broker.Client("sub")
	NewWithT(t).Expect(err).To(BeNil())

	go func() {
		err = csub.WithTopic(topic).Subscribe(func(cli mqtt.Client, msg mqtt.Message) {
			pl := &PayloadBody{}
			ts := time.Now()
			NewWithT(t).Expect(json.Unmarshal(msg.Payload(), pl)).To(BeNil())
			fmt.Printf("topic: %s cst: %dms\n", msg.Topic(), ts.UnixMilli()-pl.PubTimestamp)
		})
		NewWithT(t).Expect(err).To(BeNil())
	}()

	num := 100
	for i := 0; i < num; i++ {
		err = cpub.WithTopic(topic).WithQoS(QOS__AT_LEAST_ONCE).WithRetain(false).
			Publish(UnsafeJsonMarshal(NewPayloadBody("payload")))
		NewWithT(t).Expect(err).To(BeNil())
		time.Sleep(time.Second)
	}

	err = cpub.Unsubscribe()
	NewWithT(t).Expect(err).To(BeNil())
	err = csub.Unsubscribe()
	NewWithT(t).Expect(err).To(BeNil())
	cpub.Disconnect()
	csub.Disconnect()
}
