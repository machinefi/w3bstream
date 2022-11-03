package main

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/depends/testutil/httptransporttestutil/server/cmd/app/routes"
)

func main() {
	ht := &httptransport.HttpTransport{
		Port: 8080,
	}
	ht.SetDefault()

	kit.Run(routes.RootRouter, ht)
}
