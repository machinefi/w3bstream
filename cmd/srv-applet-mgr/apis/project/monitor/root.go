package monitor

import (
	"github.com/iotexproject/Bumblebee/kit/httptransport"
	"github.com/iotexproject/Bumblebee/kit/kit"
)

var Root = kit.NewRouter(httptransport.Group("/project/monitor"))

func init() {
	Root.Register(kit.NewRouter(&CreateMonitor{}))
}
