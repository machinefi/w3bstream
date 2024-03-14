package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/base/consts"
	base "github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/conf/postgres"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/timer"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/modules/robot_notifier"
	"github.com/machinefi/w3bstream/pkg/modules/robot_notifier/lark"
	"github.com/machinefi/w3bstream/pkg/types"
)

func main() {
	ep := base.Endpoint{}
	err := ep.UnmarshalText([]byte("postgres://w3bstream:8ShjeQUc@d';d4n@34.172.94.245:5432/w3bstream?sslmode=disable"))
	if err != nil {
		panic(err)
	}

	nc := &types.RobotNotifierConfig{
		Vendor: "lark",
		Env:    "prod",
		URL:    "https://open.larksuite.com/open-apis/bot/v2/hook/f8d7cd45-4b45-40fe-9635-5e2f85e19155",
		Secret: "vztL7BIOyDw10XEd9H5B6",
		SignFn: nil,
	}
	nc.Init()

	_ = os.Setenv(consts.EnvProjectName, "srv-pebble-pending-monitor")
	_ = os.Setenv(consts.EnvProjectVersion, "0.0.1")

	ctx := types.WithRobotNotifierConfig(context.Background(), nc)

	db := postgres.Endpoint{
		Master:   ep,
		PoolSize: 1,
		Database: models.DB,
	}
	db.SetDefault()
	if err := db.Init(); err != nil {
		panic(err)
	}

	count := int64(0)
	msg := ""
	interval := time.Minute
	threshold := int64(100)

	for {
		cost := timer.Start()
		err = db.QueryAndScan(builder.Expr("SELECT count(1) FROM applet_management.t_event WHERE f_project_id = 1456942923637714945 AND f_stage = 1"), &count)
		du := cost()
		now := time.Now().Format("2006-01-02T15:04:05")

		if err != nil {
			msg = fmt.Sprintf("[%s] query failed: %v database cost: %ds", now, err, int(du.Seconds()))
			// fmt.Println(msg)
			time.Sleep(interval)
			continue
		}
		msg = fmt.Sprintf("[%s] pebble task pending: %d query cost: %ds", now, count, int(du.Seconds()))
		fmt.Println(msg)
		if count > threshold {
			goto PUSH
		}
		time.Sleep(interval)
		continue
	PUSH:
		content, _ := lark.Build(ctx, "Pebble Task Pending", "WARNING", msg)
		if len(content) > 0 {
			_ = robot_notifier.Push(ctx, content)
		}
		time.Sleep(interval)
		continue
	}
}
