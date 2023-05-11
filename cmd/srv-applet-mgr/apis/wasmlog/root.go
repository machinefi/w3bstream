package wasmlog

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

var Root = kit.NewRouter(httptransport.Group("/wasmlog"))

func init() {
	Root.Register(kit.NewRouter(&RemoveWasmLogByInstanceID{}))
}
