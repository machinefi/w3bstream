package monitor

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

var Root = kit.NewRouter(httptransport.Group("/monitor"))

func init() {
	Root.Register(kit.NewRouter(&CreateContractLog{}))
	Root.Register(kit.NewRouter(&RemoveContractLog{}))
}
