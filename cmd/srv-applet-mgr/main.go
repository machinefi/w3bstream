package main

import (
	"context"
	"sync"
	"time"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/global"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/tasks"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/modules/account"
	"github.com/machinefi/w3bstream/pkg/modules/blockchain"
	"github.com/machinefi/w3bstream/pkg/modules/cronjob"
	"github.com/machinefi/w3bstream/pkg/modules/deploy"
	"github.com/machinefi/w3bstream/pkg/modules/operator"
	"github.com/machinefi/w3bstream/pkg/modules/project"
	"github.com/machinefi/w3bstream/pkg/types"
)

var app = global.App

func main() {
	app.AddCommand("migrate", func(args ...string) {
		global.Migrate()
	})
	app.Execute(func(args ...string) {
		BatchRun(
			func() {
				kit.Run(apis.RootMgr, global.Server())
			},
			func() {
				kit.Run(apis.RootEvent, global.EventServer())
			},
			func() {
				kit.Run(tasks.Root, global.TaskServer())
			},
			func() {
				if err := project.Init(global.Context); err != nil {
					panic(err)
				}
			},
			func() {
				if err := deploy.Init(global.Context); err != nil {
					panic(err)
				}
			},
			func() {
				l := types.MustLoggerFromContext(global.Context)

				_, l = l.Start(context.Background(), "init.CreateAdmin")
				defer l.End()

				passwd, err := account.CreateAdminIfNotExist(global.Context)
				if err != nil {
					l.Panic(err)
				}
				if passwd == "" {
					l.Info("admin already exists")
				} else {
					l.Info("admin created, default password is: '%s'", passwd)
				}
			},
			func() {
				l := types.MustLoggerFromContext(global.Context)

				_, l = l.Start(context.Background(), "init.InitChainDB")
				defer l.End()

				if err := blockchain.InitChainDB(global.Context); err != nil {
					l.Panic(err)
					return
				}
			},
			func() {
				blockchain.Monitor(global.Context)
			},
			func() {
				cronjob.Run(global.Context)
			},
			func() {
				operator.Migrate(global.Context)
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
