package configuration

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/modules/access_key"
)

var Root = kit.NewRouter(httptransport.Group("/configuration"))

func init() {
	Root.Register(kit.NewRouter(&EthClient{}))
	Root.Register(kit.NewRouter(&ChainConfig{}))

	access_key.RouterRegister(Root, enums.ApiGroupConfiguration, enums.ApiGroupConfigurationDesc)
}
