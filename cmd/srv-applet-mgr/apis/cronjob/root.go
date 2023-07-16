package cronjob

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/modules/access_key"
)

var Root = kit.NewRouter(httptransport.Group("/cronjob"))

func init() {
	Root.Register(kit.NewRouter(&CreateCronJob{}))
	Root.Register(kit.NewRouter(&RemoveCronJob{}))
	Root.Register(kit.NewRouter(&ListCronJob{}))

	access_key.RouterRegister(Root, enums.ApiGroupCronjob, enums.ApiGroupCronjobDesc)
}
