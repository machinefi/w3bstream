package publisher

import (
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/modules/access_key"
)

var Root = kit.NewRouter(httptransport.Group("/publisher"))

func init() {
	Root.Register(kit.NewRouter(&GetPublisher{}))
	Root.Register(kit.NewRouter(&middleware.ProjectProvider{}, &ListPublisher{}))
	Root.Register(kit.NewRouter(&RemovePublisher{}))
	Root.Register(kit.NewRouter(&middleware.ProjectProvider{}, &BatchRemovePublisher{}))
	Root.Register(kit.NewRouter(&middleware.ProjectProvider{}, &CreatePublisher{}))
	Root.Register(kit.NewRouter(&middleware.ProjectProvider{}, &UpdatePublisher{}))

	access_key.RouterRegister(Root, enums.ApiGroupPublisher, enums.ApiGroupPublisherDesc)
}
