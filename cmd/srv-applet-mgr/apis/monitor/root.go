package monitor

import (
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/modules/access_key"
)

var Root = kit.NewRouter(httptransport.Group("/monitor"))

func init() {
	Root.Register(kit.NewRouter(&middleware.ProjectProvider{}, &CreateContractLog{}))
	Root.Register(kit.NewRouter(&middleware.ProjectProvider{}, &CreateChainTx{}))
	Root.Register(kit.NewRouter(&middleware.ProjectProvider{}, &CreateChainHeight{}))
	Root.Register(kit.NewRouter(&middleware.ProjectProvider{}, &RemoveContractLog{}))
	Root.Register(kit.NewRouter(&middleware.ProjectProvider{}, &RemoveChainTx{}))
	Root.Register(kit.NewRouter(&middleware.ProjectProvider{}, &RemoveChainHeight{}))
	Root.Register(kit.NewRouter(&middleware.ProjectProvider{}, &ControlContractLog{}))
	Root.Register(kit.NewRouter(&middleware.ProjectProvider{}, &ControlChainTx{}))
	Root.Register(kit.NewRouter(&middleware.ProjectProvider{}, &ControlChainHeight{}))

	access_key.RouterRegister(Root, enums.ApiGroupMonitor, enums.ApiGroupMonitorDesc)
}
