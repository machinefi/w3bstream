package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/iotexproject/Bumblebee/kit/kit"

	"github.com/iotexproject/w3bstream/cmd/srv-applet-mgr/apis"
	"github.com/iotexproject/w3bstream/cmd/srv-applet-mgr/global"
	"github.com/iotexproject/w3bstream/pkg/modules/account"
	"github.com/iotexproject/w3bstream/pkg/modules/blockchain"
	"github.com/iotexproject/w3bstream/pkg/modules/deploy"
)

var app = global.App

func main() {
	app.AddCommand("migrate", func(args ...string) {
		global.Migrate()
	})

	app.Execute(func(args ...string) {
		BatchRun(
			func() {
				kit.Run(apis.Root, global.Server())
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

				passwd, err := account.CreateAdminIfNotExist(ctx)
				if err != nil {
					log.Panicf("create admin faild: %s", err.Error())
					return
				}
				if passwd == "" {
					log.Printf("admin already created")
					return
				}
				log.Printf("admin created, default password: iotex.W3B.admin")
			},
			func() {
				go blockchain.ListenContractlog(global.WithContext(context.Background()))
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
