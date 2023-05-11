package resource

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

var Root = kit.NewRouter(httptransport.Group("/resource"))

func init() {
	Root.Register(kit.NewRouter(&ListResources{}))
	Root.Register(kit.NewRouter(&RemoveResource{}))
	Root.Register(kit.NewRouter(&DownloadResource{}))
}
