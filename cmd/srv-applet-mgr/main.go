package main

import (
	"context"
	"sync"
	"time"

	"github.com/iotexproject/Bumblebee/kit/kit"
	"github.com/iotexproject/w3bstream/pkg/types"

	"github.com/iotexproject/w3bstream/cmd/srv-applet-mgr/apis"
	"github.com/iotexproject/w3bstream/cmd/srv-applet-mgr/global"
	"github.com/iotexproject/w3bstream/pkg/depends/protocol/eventpb"
	"github.com/iotexproject/w3bstream/pkg/modules/account"
	"github.com/iotexproject/w3bstream/pkg/modules/blockchain"
	"github.com/iotexproject/w3bstream/pkg/modules/deploy"
	"github.com/iotexproject/w3bstream/pkg/modules/event"
	"github.com/iotexproject/w3bstream/pkg/modules/project"
)

var app = global.App

func main() {
	// app.AddCommand("migrate", func(args ...string) {
	// 	global.Migrate()
	// })
	app.Execute(func(args ...string) {
		BatchRun(
			global.Migrate,
			func() {
				kit.Run(apis.Root, global.Server())
			},
			func() {
				if err := project.InitChannels(
					global.WithContext(context.Background()),
					func(ctx context.Context, channel string, data *eventpb.Event) (interface{}, error) {
						return event.OnEventReceived(ctx, channel, data)
					},
				); err != nil {
					panic(err)
				}
			},
			func() {
				if err := deploy.StartInstances(
					global.WithContext(context.Background()),
				); err != nil {
					panic(err)
				}
			},
			func() {
				ctx := global.WithContext(context.Background())
				l := types.MustLoggerFromContext(ctx)

				_, l = l.Start(ctx, "init.CreateAdmin")
				defer l.End()

				passwd, err := account.CreateAdminIfNotExist(ctx)
				if err != nil {
					l.Panic(err)
					return
				}
				if passwd == "" {
					l.Info("admin already created")
					return
				}
			},
			func() {
				ctx := global.WithContext(context.Background())
				l := types.MustLoggerFromContext(ctx)

				_, l = l.Start(ctx, "init.InitChainDB")
				defer l.End()

				if err := blockchain.InitChainDB(ctx); err != nil {
					l.Panic(err)
					return
				}
			},
			func() {
				blockchain.Monitor(global.WithContext(context.Background()))
			},
		)
	})
}

func BatchRun(commands ...func()) {
	wg := &sync.WaitGroup{}

	for i := range commands {
		cmd := commands[i]
		wg.Add(1)

		go func() {
			defer wg.Done()
			cmd()
			time.Sleep(200 * time.Millisecond)
		}()
	}
	wg.Wait()
}
