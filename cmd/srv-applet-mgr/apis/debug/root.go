package debug

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

var Root = kit.NewRouter(httptransport.Group("/debug"))

func init() {
	Root.Register(kit.NewRouter(&FetchInstances{}))
	Root.Register(kit.NewRouter(&GetInstance{}))
}
