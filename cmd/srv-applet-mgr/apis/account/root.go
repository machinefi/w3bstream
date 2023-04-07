package account

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

var (
	Root         = kit.NewRouter(httptransport.Group("/account"))
	RegisterRoot = kit.NewRouter(httptransport.Group("/register"))
)

func init() {
	Root.Register(kit.NewRouter(&UpdatePasswordByAccountID{}))
	Root.Register(kit.NewRouter(&GetOperatorAddr{}))

	RegisterRoot.Register(kit.NewRouter(&CreateAccountByUsernameAndPassword{}))
}
