package resource

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/modules/access_key"
)

var Root = kit.NewRouter(httptransport.Group("/resource"))

func init() {
	Root.Register(kit.NewRouter(&ListResources{}))
	Root.Register(kit.NewRouter(&RemoveResource{}))
	Root.Register(kit.NewRouter(&DownloadResource{}))
	Root.Register(kit.NewRouter(&GetDownloadResourceUrl{}))

	access_key.RouterRegister(Root, enums.ApiGroupResource, enums.ApiGroupResourceDesc)
}
