package deploy

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

var Root = kit.NewRouter(httptransport.Group("/deploy"))

func init() {
	Root.Register(kit.NewRouter(&CreateInstance{}))
	Root.Register(kit.NewRouter(&GetInstanceByInstanceID{}))
	Root.Register(kit.NewRouter(&GetInstanceByAppletID{}))
	Root.Register(kit.NewRouter(&ControlInstance{}))
	Root.Register(kit.NewRouter(&ReDeployInstance{}))
}
