package project_config

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

var Root = kit.NewRouter(httptransport.Group("/project_config"))

func init() {
	Root.Register(kit.NewRouter(&GetProjectSchema{}))
	Root.Register(kit.NewRouter(&GetProjectEnv{}))
	Root.Register(kit.NewRouter(&CreateProjectSchema{}))
	Root.Register(kit.NewRouter(&CreateOrUpdateProjectEnv{}))
}
