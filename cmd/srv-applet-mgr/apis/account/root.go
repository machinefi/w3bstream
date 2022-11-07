package account

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

var Root = kit.NewRouter(httptransport.Group("/account"))

func init() {
	Root.Register(kit.NewRouter(&CreateAccount{}))
	Root.Register(kit.NewRouter(&UpdatePasswordByAccountID{}))
}
