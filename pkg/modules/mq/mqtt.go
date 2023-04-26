package mq

import (
	"context"
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	confmqtt "github.com/machinefi/w3bstream/pkg/depends/conf/mqtt"
	"github.com/machinefi/w3bstream/pkg/depends/protocol/eventpb"
	"github.com/machinefi/w3bstream/pkg/depends/x/mapx"
	"github.com/machinefi/w3bstream/pkg/types"
)

type OnMessage func(ctx context.Context, channel string, data *eventpb.Event) (interface{}, error)

var channels = mapx.New[string, *ChannelContext]()

type ChannelContext struct {
	ctx    context.Context
	cancel context.CancelFunc
	Name   string
	cli    *confmqtt.Client
	hdl    OnMessage
}

func (cc *ChannelContext) Run(ctx context.Context) {
	l := types.MustLoggerFromContext(ctx)

	_, _l := l.Start(ctx, "ChannelContext.Run")
	defer _l.End()

	mqHandler := func(cli mqtt.Client, msg mqtt.Message) {
		_, l := l.Start(cc.ctx, "OnMessage:"+cc.Name)
		defer l.End()

		pl := msg.Payload()
		ev := &eventpb.Event{}
		err := json.Unmarshal(pl, ev)
		if err != nil {
			l.Error(err)
			return
		}
		// trans event to event proxy and call c.hdl(ctx, payload)
		_, err = cc.hdl(ctx, cc.Name, ev)
		if err != nil {
			l.Error(err)
		}
		l.WithValues("payload", ev).Info("sub handled")
	}
	if err := cc.cli.Subscribe(mqHandler); err != nil {
		return
	}
	defer func() {
		_l.Info("channel closed")
		if err := cc.cli.Unsubscribe(); err != nil {
			l.Error(err)
		}
	}()

	<-cc.ctx.Done()
}

func (cc *ChannelContext) Stop() { cc.cancel() }

func CreateChannel(ctx context.Context, prjName string, hdl OnMessage) error {
	l := types.MustLoggerFromContext(ctx)
	defer l.End()

	_, l = l.Start(ctx, "CreateChannel")
	defer l.End()

	l = l.WithValues("project_name", prjName)

	broker := types.MustMqttBrokerFromContext(ctx)

	cli, err := broker.Client(prjName)
	if err != nil {
		l.Error(err)
		return err
	}

	cctx := &ChannelContext{
		Name: prjName,
		cli:  cli.WithTopic(prjName),
		hdl:  hdl,
	}
	cctx.ctx, cctx.cancel = context.WithCancel(context.Background())
	channels.Store(prjName, cctx)

	go cctx.Run(ctx)

	l.Info("channel started")
	return nil
}

func StopChannel(prjName string) {
	c, ok := channels.LoadAndRemove(prjName)
	if !ok {
		return
	}
	c.Stop()
}
