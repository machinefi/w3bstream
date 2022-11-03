package main

import (
	"flag"
	"fmt"

	confapp "github.com/machinefi/w3bstream/pkg/depends/conf/app"
	conflog "github.com/machinefi/w3bstream/pkg/depends/conf/log"
	confmqtt "github.com/machinefi/w3bstream/pkg/depends/conf/mqtt"
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

	flag.StringVar(&id, "i", "", "publish client id")
	flag.StringVar(&topic, "t", "", "publish topic")
	flag.StringVar(&payload, "c", "", "publish content")
	flag.Parse()
}

var (
	id      string
	topic   string
	payload string
)

func main() {
	c, err := mqtt.Client(id)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = c.WithTopic(topic).Publish(payload)
	if err != nil {
		fmt.Println(err)
		return
	}
}
