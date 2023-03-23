package login

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

var Root = kit.NewRouter(httptransport.Group("/login"))

func init() {
	Root.Register(kit.NewRouter(&LoginByUsername{}))
	Root.Register(kit.NewRouter(&LoginByEthAddress{}))
}
