package account_access

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

var (
	Root = kit.NewRouter(httptransport.Group("/access_key"))
)

func init() {
	Root.Register(kit.NewRouter(&CreateAccountAccessKey{}))
	Root.Register(kit.NewRouter(&DeleteAccountAccessKeyByName{}))
}
