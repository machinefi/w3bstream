package mq

import (
	"context"
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	conflog "github.com/machinefi/w3bstream/pkg/depends/conf/log"
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
	_, l := conflog.FromContext(ctx).Start(ctx)
	defer l.End()

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
		_, err = cc.hdl(ctx, cc.Name, ev)
		if err != nil {
			l.Error(err)
		}
		l.WithValues("payload", ev).Info("sub handled")
	}
	l = l.WithValues("cid", cc.cli.Cid(), "topic", cc.Name)
	if err := cc.cli.Subscribe(mqHandler); err != nil {
		l.Error(err)
		return
	}
	l.Info("start subscribing")
	defer func() {
		l.Info("stop subscribing and mq client closed")
		if err := cc.cli.Unsubscribe(); err != nil {
			l.Error(err)
		}
		cc.cli.Disconnect()
	}()

	<-cc.ctx.Done()
}

func (cc *ChannelContext) Stop() { cc.cancel() }

func CreateChannel(ctx context.Context, prjName string, hdl OnMessage) error {
	l := types.MustLoggerFromContext(ctx)
	defer l.End()

	_, l = l.Start(ctx)
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

	return nil
}

func StopChannel(channel string) {
	c, ok := channels.LoadAndRemove(channel)
	if !ok {
		return
	}
	c.Stop()
}
