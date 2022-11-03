package tasks

import "github.com/machinefi/w3bstream/pkg/depends/kit/kit"

var Root = kit.NewRouter()

func init() {
	Root.Register(kit.NewRouter(&HandleEvent{}))
}
