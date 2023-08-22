package main

import (
	"sync"

	"github.com/machinefi/w3bstream/cmd/pub_bench/global"
	"github.com/machinefi/w3bstream/cmd/pub_bench/types"
)

func main() {
	ctx := global.Context
	chs := types.MustChannelsFromContext(ctx)

	wg := &sync.WaitGroup{}
	for _, c := range chs {
		if c.IsZero() {
			continue
		}
		if err := c.Subscribe(ctx); err != nil {
			panic(err)
		}
		wg.Add(1)
		go func(c *types.Channel) {
			c.StartPublish(ctx)
			wg.Done()
		}(c)
	}
	wg.Wait()
}
