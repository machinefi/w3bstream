package login

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/modules/access_key"
)

var Root = kit.NewRouter(httptransport.Group("/login"))

func init() {
	Root.Register(kit.NewRouter(&LoginByUsername{}))
	Root.Register(kit.NewRouter(&LoginByEthAddress{}))

	access_key.RouterRegister(Root, enums.ApiGroupLogin, enums.ApiGroupLoginDesc)
}
