package publisher

import (
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

var Root = kit.NewRouter(httptransport.Group("/publisher"))

func init() {
	Root.Register(kit.NewRouter(&middleware.ProjectProvider{}, &ListPublisher{}))
	Root.Register(kit.NewRouter(&middleware.ProjectProvider{}, &CreatePublisher{}))
	Root.Register(kit.NewRouter(&middleware.ProjectProvider{}, &UpdatePublisher{}))
	Root.Register(kit.NewRouter(&RemovePublisher{}))
}
