package project_config

import (
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

var Root = kit.NewRouter(httptransport.Group("/project_config"))

func init() {
	Root.Register(kit.NewRouter(&middleware.ProjectProvider{}, &GetProjectSchema{}))
	Root.Register(kit.NewRouter(&middleware.ProjectProvider{}, &GetProjectEnv{}))
	Root.Register(kit.NewRouter(&middleware.ProjectProvider{}, &CreateProjectSchema{}))
	Root.Register(kit.NewRouter(&middleware.ProjectProvider{}, &CreateOrUpdateProjectEnv{}))
}
