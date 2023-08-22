package tasks

import "github.com/machinefi/w3bstream/pkg/depends/kit/kit"

var Root = kit.NewRouter()

func init() {
	Root.Register(kit.NewRouter(&HandleEvent{}))
	Root.Register(kit.NewRouter(&DbLogStoring{}))
	Root.Register(kit.NewRouter(&EventLog{}))
	Root.Register(kit.NewRouter(&EventLogCleanup{}))
}
