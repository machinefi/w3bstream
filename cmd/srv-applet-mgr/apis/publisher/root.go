package publisher

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

var Root = kit.NewRouter(httptransport.Group("/publisher"))

func init() {
	Root.Register(kit.NewRouter(&ListPublisher{}))
	Root.Register(kit.NewRouter(&CreatePublisher{}))
	Root.Register(kit.NewRouter(&UpdatePublisher{}))
	Root.Register(kit.NewRouter(&RemovePublisher{}))
}
