package tag

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

var Root = kit.NewRouter(httptransport.Group("/tag"))

func init() {
	Root.Register(kit.NewRouter(&CreateTag{}))
	Root.Register(kit.NewRouter(&ListTag{}))
	Root.Register(kit.NewRouter(&RemoveTag{}))
}
