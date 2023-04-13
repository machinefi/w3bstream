package main

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransportgen/testdata/server/cmd/app/routes"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

func main() {
	ht := &httptransport.HttpTransport{
		Port: 8080,
	}
	ht.SetDefault()

	kit.Run(routes.RootRouter, ht)
}
