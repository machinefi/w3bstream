package monitor

import (
	"github.com/iotexproject/Bumblebee/kit/httptransport"
	"github.com/iotexproject/Bumblebee/kit/kit"
)

var Root = kit.NewRouter(httptransport.Group("/monitor/contractlog"))

func init() {
	Root.Register(kit.NewRouter(&CreateContractlog{}))
}
