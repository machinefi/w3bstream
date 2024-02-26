package main

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/global"
	"github.com/machinefi/w3bstream/pkg/depends/conf/logger"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/depends/kit/logr"
	"github.com/machinefi/w3bstream/pkg/modules/account"
	"github.com/machinefi/w3bstream/pkg/modules/blockchain"
	"github.com/machinefi/w3bstream/pkg/modules/cronjob"
	"github.com/machinefi/w3bstream/pkg/modules/deploy"
	"github.com/machinefi/w3bstream/pkg/modules/event"
	"github.com/machinefi/w3bstream/pkg/modules/metrics"
	"github.com/machinefi/w3bstream/pkg/modules/operator"
	"github.com/machinefi/w3bstream/pkg/modules/project"
	"github.com/machinefi/w3bstream/pkg/modules/robot_notifier"
	"github.com/machinefi/w3bstream/pkg/modules/robot_notifier/lark"
	"github.com/machinefi/w3bstream/pkg/modules/trafficlimit"
	"github.com/machinefi/w3bstream/pkg/types"
)

var app = global.App

func init() {
	global.Migrate()
}

func main() {
	ctx, l := logger.NewSpanContext(global.WithContext(context.Background()), "main")
	defer l.End()

	var (
		sigProjectsInitialized = make(chan struct{})
		projectIDs             []types.SFID
	)
	app.Execute(func(args ...string) {
		BatchRun(
			func() {
				kit.Run(apis.RootMgr, global.Server())
			},
			func() {
				kit.Run(apis.RootEvent, global.EventServer())
			},
			func() {
				kit.Run(apis.RootDebug, global.DebugServer())
			},
			func() {
				ctx, l := logr.Start(ctx, "main.InitProjects")
				defer l.End()

				passwd, err := account.CreateAdminIfNotExist(ctx)
				if err != nil {
					l.Error(err)
					panic(err)
				}
				if passwd == "" {
					l.Info("admin already exists")
				} else {
					l.Info("admin created, default password is: '%s'", passwd)
				}

				if err := deploy.Init(ctx); err != nil {
					l.Error(err)
					panic(err)
				}
				if projectIDs, err = project.Init(ctx); err != nil {
					l.Error(err)
					panic(err)
				}
				l.Info("all projects initialized")
				sigProjectsInitialized <- struct{}{}
			},
			func() {
				if err := trafficlimit.Init(ctx); err != nil {
					panic(err)
				}
			},
			func() {
				if err := blockchain.InitChainDB(ctx); err != nil {
					l.Error(err)
					panic(err)
				}
			},
			func() {
				blockchain.Monitor(ctx)
			},
			func() {
				cronjob.Run(ctx)
			},
			func() {
				operator.Migrate(ctx)
			},
			func() {
				metrics.Init(ctx)
			},
			func() {
				<-sigProjectsInitialized
				wl, _ := types.ProjectWhiteListFromContext(ctx)
				bl, _ := types.ProjectBlackListFromContext(ctx)
				projects := make([]string, 0)
				for _, prj := range projectIDs {
					if slices.Contains(wl, prj) {
						sche := event.NewEventHandleScheduler(time.Minute/2, 300, prj)
						go sche.Run(ctx)
						projects = append(projects, prj.String())
						continue
					}
					if !slices.Contains(bl, prj) {
						sche := event.NewEventHandleScheduler(time.Minute/2, 300, prj)
						projects = append(projects, prj.String())
						go sche.Run(ctx)
					}
				}

				body, err := lark.Build(ctx, "Projects Processes", "INFO", strings.Join(projects, "\n"))
				if err != nil {
					return
				}
				_ = robot_notifier.Push(ctx, body, nil)
			},
			func() {
				sche := event.NewEventCleanupScheduler(time.Hour, 14*time.Hour*24)
				sche.Run(ctx)
			},
			func() {
				body, err := lark.Build(
					ctx,
					"service started",
					"INFO",
					fmt.Sprintf("service started at: %s", types.Timestamp{Time: time.Now()}.String()),
				)
				if err != nil {
					return
				}
				_ = robot_notifier.Push(ctx, body, nil)
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
