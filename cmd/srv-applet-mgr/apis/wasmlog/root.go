package wasmlog

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/modules/access_key"
)

var Root = kit.NewRouter(httptransport.Group("/wasmlog"))

func init() {
	Root.Register(kit.NewRouter(&RemoveWasmLogByInstanceID{}))

	access_key.RouterRegister(Root, enums.ApiGroupWasmLog, enums.ApiGroupWasmLogDesc)
}
