package main

import (
	"sync"
	"time"

	"github.com/machinefi/w3bstream/cmd/srv-gw-device/apis"
	"github.com/machinefi/w3bstream/cmd/srv-gw-device/global"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
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
