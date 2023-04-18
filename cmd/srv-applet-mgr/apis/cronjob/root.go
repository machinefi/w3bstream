package cronjob

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

var Root = kit.NewRouter(httptransport.Group("/cronjob"))

func init() {
	Root.Register(kit.NewRouter(&CreateCronJob{}))
	Root.Register(kit.NewRouter(&RemoveCronJob{}))
}
