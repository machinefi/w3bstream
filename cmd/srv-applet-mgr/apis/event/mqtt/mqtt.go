package mqtt

import (
	"context"
	"errors"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/iotexproject/w3bstream/pkg/modules/event/proxy"
	"github.com/iotexproject/w3bstream/pkg/types"
)

func Run(ctx context.Context, username, password, broker string, port int) {
	logger := types.MustLoggerFromContext(ctx)

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetUsername(username)
	opts.SetPassword(password)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		logger.Fatal(token.Error())
	}

	sub(ctx, client)
}

func sub(ctx context.Context, client mqtt.Client) {
	logger := types.MustLoggerFromContext(ctx)
	topic := "event/#" // mqtt root topic, TODO move to config

	token := client.Subscribe(topic, 1, func(client mqtt.Client, msg mqtt.Message) {
		// TODO validate publisherID
		var publisherID, projectID, appletID, handler string
		n, err := fmt.Sscanf(msg.Topic(), "event/%s/%s/%s/%s", &publisherID, &projectID, &appletID, &handler)
		if err != nil {
			logger.Error(err)
			return
		}
		if n != 4 {
			logger.Error(errors.New("invalid topic schema"))
			return
		}

		proxy.Proxy(ctx, &event{
			handler:     handler,
			projectID:   projectID,
			appletID:    appletID,
			publisherID: publisherID,
			data:        msg.Payload(),
		})

	})
	token.Wait()
}
