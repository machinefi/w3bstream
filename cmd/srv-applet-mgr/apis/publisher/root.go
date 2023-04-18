package publisher

import (
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

var (
	Root     = kit.NewRouter(httptransport.Group("/publisher"))
	Provider = kit.NewRouter(httptransport.Group("/publisher/x"))
)

func init() {
	Provider.Register(kit.NewRouter(&middleware.ProjectProvider{}, &ListPublisher{}))
	Provider.Register(kit.NewRouter(&middleware.ProjectProvider{}, &ListPublisherDetail{}))
	Provider.Register(kit.NewRouter(&middleware.ProjectProvider{}, &CreatePublisher{}))
	Root.Register(kit.NewRouter(&UpdatePublisher{}))
	Root.Register(kit.NewRouter(&RemovePublisher{}))
	Provider.Register(kit.NewRouter(&middleware.ProjectProvider{}, &BatchRemoveByPublisherIDs{}))
}
