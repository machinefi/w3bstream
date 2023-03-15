package login

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

var LoginRoot = kit.NewRouter(httptransport.Group("/login"))
var NonceRoot = kit.NewRouter(httptransport.Group("/nonce"))

func init() {
	LoginRoot.Register(kit.NewRouter(&LoginByUsername{}))
	LoginRoot.Register(kit.NewRouter(&LoginByEthAddress{}))
	NonceRoot.Register(kit.NewRouter(&GetNonceByEthAddress{}))
}
