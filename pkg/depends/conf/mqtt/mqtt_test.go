package mqtt_test

import (
	"testing"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	. "github.com/machinefi/w3bstream/pkg/depends/conf/mqtt"
)

func TestBroker(t *testing.T) {
	topic := "test_demo"
	server := types.Endpoint{}

	err := server.UnmarshalText([]byte("mqtt://broker.emqx.io:1883"))
	NewWithT(t).Expect(err).To(BeNil())

	broker := &Broker{Server: server}
	broker.SetDefault()
	broker.Init()

	c1, err := broker.Client("c1")
	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(c1).NotTo(BeNil())

	c2, err := broker.Client("c2")
	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(c2).NotTo(BeNil())

	err = c1.WithTopic(topic).WithQoS(QOS__AT_LEAST_ONCE).WithRetain(false).
		Publish("testpublish")
	NewWithT(t).Expect(err).To(BeNil())

	c2.WithTopic(topic).Subscribe(func(c mqtt.Client, msg mqtt.Message) {
		NewWithT(t).Expect(string(msg.Payload())).To(Equal("testpublish"))
	})
}
